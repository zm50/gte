package gcore

import (
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/core"
	"github.com/zm50/gte/gconf"
	"github.com/zm50/gte/glog"
	"github.com/zm50/gte/trait"
)

// ConnMgr 连接管理模块，管理客户端的连接，监听待读取数据的连接，并将连接提交给下游的消息分发模块进行处理
type ConnMgr[T any] struct {
	epfd int

	timeout int

	events []syscall.EpollEvent

	dispatcher trait.Dispatcher[T]

	// key: fd, value: Conn
	connShards *core.KVShards[int32, trait.Connection[T]]

	connStartHook func(conn trait.Connection[T])

	connStopHook func(conn trait.Connection[T])

	connNotActiveHook func(conn trait.Connection[T])

	keepAliveMgr trait.KeepAliveMgr[T]

	connSignalQueue []chan trait.ConnSignal[T]

	wg *sync.WaitGroup
}

var _ trait.ConnMgr[int] = (*ConnMgr[int])(nil)

// NewConnMgr 新建一个连接管理的实例
func NewConnMgr[T any](timeout int, eventSize int, taskMgr trait.TaskMgr[T]) (*ConnMgr[T], error) {
	// 创建一个epoll句柄
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	connSignalQueues := make([]chan trait.ConnSignal[T], gconf.Config.ConnSignalQueues())
	for i := 0; i < len(connSignalQueues); i++ {
		connSignalQueues[i] = make(chan trait.ConnSignal[T], gconf.Config.ConnSignalQueueLen())
	}

	connShards := core.NewKVShards[int32, trait.Connection[T]](gconf.Config.ConnShardCount())

	// 创建一个连接管理器
	connMgr := &ConnMgr[T]{
		epfd:            epfd,
		timeout:         timeout,
		events:          make([]syscall.EpollEvent, eventSize),
		connShards:      connShards,
		connSignalQueue: connSignalQueues,
		wg:              &sync.WaitGroup{},
	}

	connMgr.dispatcher = NewDispatcher(connMgr, taskMgr)

	connMgr.keepAliveMgr = NewKeepAliveMgr(connMgr, connShards.Shards())

	return connMgr, nil
}

// Get 在连接管理器中查群连接
func (e *ConnMgr[T]) Get(fd int32) (trait.Connection[T], bool) {
	return e.connShards.Get(fd)
}

// Add 在连接管理器中添加连接
func (e *ConnMgr[T]) Add(conn trait.Connection[T]) error {
	sock, err := conn.File()
	if err != nil {
		glog.Error("get socket file descriptor error:", err)
		return err
	}
	fd := sock.Fd()

	if _, ok := e.Get(int32(fd)); ok {
		glog.Error("connection already exists, conn fd:", fd)
		return errors.Errorf("connection already exists, conn fd: %d", fd)
	}

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLRDHUP,
		Fd:     int32(fd),
	}
	err = syscall.EpollCtl(e.epfd, syscall.EPOLL_CTL_ADD, int(fd), &event)
	if err != nil {
		glog.Error("epoll ctl add error:", err)
		return err
	}

	e.connShards.Set(int32(fd), conn)

	// 通知连接信号处理队列
	e.PushConnSignal(NewConnSignal[T](conn, constant.ConnStartSignal))

	return nil
}

// Del 在连接管理器中删除连接
func (e *ConnMgr[T]) Del(fd int32) error {
	// 通知连接信号处理队列
	conn, ok := e.Get(fd)
	if ok {
		e.PushConnSignal(NewConnSignal[T](conn, constant.ConnStopSignal))
	} else {
		glog.Error("call conn stop hook failed, connection not found, conn fd:", fd)
	}

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}
	err := syscall.EpollCtl(e.epfd, syscall.EPOLL_CTL_DEL, int(fd), &event)
	if err != nil {
		glog.Error("epoll ctl del error:", err)
		return err
	}

	e.connShards.Del(fd)

	return nil
}

// Wait 等待事件发生
func (e *ConnMgr[T]) Wait() (int, error) {
	n, err := syscall.EpollWait(e.epfd, e.events, e.timeout)

	return n, err
}

// BatchCommit 批量处理待读取的连接
func (e *ConnMgr[T]) BatchCommit(n int) {
	e.wg.Add(n)

	for n > 0 {
		n--
		event := e.events[n]
		fd := event.Fd

		if event.Events&syscall.EPOLLRDHUP != 0 {
			// 连接关闭事件处理
			e.Del(event.Fd)
			continue
		}

		conn, ok := e.Get(fd)
		if !ok {
			glog.Error("connection not found, conn fd:", fd)
			continue
		}

		// 提交待读取的连接
		e.dispatcher.Commit(conn)
	}

	e.wg.Wait()
}

// Start 启动连接管理器
func (e *ConnMgr[T]) Start() {
	e.dispatcher.Start()

	e.keepAliveMgr.Start()

	e.StartConnSignalHookWorkers()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer syscall.Close(e.epfd)

	delay := time.Duration(gconf.Config.EpollTimeout()) * time.Millisecond

	for {
		n, err := e.Wait()
		if err != nil {
			glog.Error("epoll wait error:", err)
			continue
		}

		if n == 0 {
			time.Sleep(delay)
			continue
		}

		now := time.Now()
		glog.Info("epoll wait events:", n)
		e.BatchCommit(n)
		glog.Infof("epoll events processed, cost: %s\n", time.Now().Sub(now).String())
	}
}

// Stop 结束连接管理器
func (e *ConnMgr[T]) Stop() {
	n := e.connShards.Count()
	for conn := range e.connShards.ValuesIter(n) {
		conn.Close()
	}

	syscall.Close(e.epfd)
}

// StartConnSignalHookWorkers 启动连接信号钩子消费者工作池
func (e *ConnMgr[T]) StartConnSignalHookWorkers() {
	for i := 0; i < len(e.connSignalQueue); i++ {
		for j := 0; j < gconf.Config.WorkersPerConnSignalQueue(); j++ {
			go e.StartConnSignalHookWorker(e.connSignalQueue[i])
		}
	}
}

// StartConnSignalHookWorker 启动连接信号钩子消费者
func (m *ConnMgr[T]) StartConnSignalHookWorker(connSignalQueue <-chan trait.ConnSignal[T]) {
	for conn := range connSignalQueue {
		switch conn.Signal() {
		case constant.ConnStartSignal:
			if m.connStartHook != nil {
				m.connStartHook(conn)
			}
		case constant.ConnStopSignal:
			if m.connStopHook != nil {
				m.connStopHook(conn)
			}
		case constant.ConnNotActiveSignal:
			if m.connNotActiveHook != nil {
				m.connNotActiveHook(conn)
			}
		default:
			glog.Error("unknown conn signal:", conn.Signal())
		}
	}
}

// OnConnStart 注册连接建立触发的钩子回调
func (m *ConnMgr[T]) OnConnStart(fn func(conn trait.Connection[T])) {
	m.connStartHook = fn
}

// OnConnStop 注册连接断开触发的钩子回调
func (m *ConnMgr[T]) OnConnStop(fn func(conn trait.Connection[T])) {
	m.connStopHook = fn
}

// OnConnNotActive 注册连接不活跃触发的钩子回调
func (m *ConnMgr[T]) OnConnNotActive(fn func(conn trait.Connection[T])) {
	m.connNotActiveHook = fn
}

// ChooseConnSignalQueue 选择连接信号处理队列
func (m *ConnMgr[T]) ChooseConnSignalQueue(connID uint64) chan<- trait.ConnSignal[T] {
	return m.connSignalQueue[connID%uint64(len(m.connSignalQueue))]
}

// PushConnSignal 推送连接信号
func (m *ConnMgr[T]) PushConnSignal(signal trait.ConnSignal[T]) {
	m.ChooseConnSignalQueue(signal.ID()) <- signal
}

// WaitGroup 等待组
func (m *ConnMgr[T]) WaitGroup() *sync.WaitGroup {
	return m.wg
}
