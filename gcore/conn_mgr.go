package gcore

import (
	"runtime"
	"syscall"
	"time"

	"github.com/go75/gte/constant"
	"github.com/go75/gte/core"
	"github.com/go75/gte/gconf"
	"github.com/go75/gte/glog"
	"github.com/go75/gte/trait"
	"github.com/pkg/errors"
)

// ConnMgr 连接管理模块，管理客户端的连接，监听待读取数据的连接，并将连接提交给下游的消息分发模块进行处理
type ConnMgr struct {
	epfd int

	timeout int

	events []syscall.EpollEvent

	readTimeout time.Duration
	maxReadTimeout time.Duration

	readDeadline time.Time
	maxReadDeadline time.Time

	dispatcher trait.Dispatcher

	// key: fd, value: Conn
	connShards *core.KVShards[int32, trait.Connection]

	connStartHook func(conn trait.Connection)

	connStopHook func(conn trait.Connection)

	connNotActiveHook func(conn trait.Connection)

	keepAliveMgr trait.KeepAliveMgr

	connSignalQueue  []chan trait.ConnSignal
}

var _ trait.ConnMgr = (*ConnMgr)(nil)

// NewConnMgr 新建一个连接管理的实例
func NewConnMgr(timeout int, eventSize int, taskMgr trait.TaskMgr) (*ConnMgr, error) {
	// 创建一个epoll句柄
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	readTimeout := time.Duration(gconf.Config.ReadTimeout()) * time.Millisecond
	maxReadTimeout := time.Duration(gconf.Config.MaxReadTimeout()) * time.Millisecond

	connSignalQueues := make([]chan trait.ConnSignal, gconf.Config.ConnSignalQueues())
	for i := 0; i < len(connSignalQueues); i++ {
		connSignalQueues[i] = make(chan trait.ConnSignal, gconf.Config.ConnSignalQueueLen())
	}

	connShards := core.NewKVShards[int32, trait.Connection](gconf.Config.ConnShardCount())

	// 创建一个连接管理器
	connMgr := &ConnMgr{
		epfd: epfd,
		timeout: timeout,
		events: make([]syscall.EpollEvent, eventSize),
		connShards: connShards,
		readTimeout: readTimeout,
		maxReadTimeout: maxReadTimeout,
		connSignalQueue: connSignalQueues,
	}

	connMgr.dispatcher = NewDispatcher(connMgr, taskMgr)

	connMgr.keepAliveMgr = NewKeepAliveMgr(connMgr, connShards.Shards())

	return connMgr, nil
}

// Get 在连接管理器中查群连接
func (e *ConnMgr) Get(fd int32) (trait.Connection, bool) {
	return e.connShards.Get(fd)
}

// Add 在连接管理器中添加连接
func (e *ConnMgr) Add(conn trait.Connection) error {
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
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}
	err = syscall.EpollCtl(e.epfd, syscall.EPOLL_CTL_ADD | syscall.EPOLLRDHUP, int(fd), &event)
	if err != nil {
		glog.Error("epoll ctl add error:", err)
		return err
	}

	e.connShards.Set(int32(fd), conn)

	// 通知连接信号处理队列
	e.PushConnSignal(NewConnSignal(conn, constant.ConnStartSignal))

	return nil
}

// Del 在连接管理器中删除连接
func (e *ConnMgr) Del(fd int32) error {
	// 通知连接信号处理队列
	conn, ok := e.Get(fd)
	if ok {
		e.PushConnSignal(NewConnSignal(conn, constant.ConnStopSignal))
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
func (e *ConnMgr) Wait() (int, error) {
	n, err := syscall.EpollWait(e.epfd, e.events, e.timeout)

	return n, err
}

// BatchCommit 批量处理待读取的连接
func (e *ConnMgr) BatchCommit(n int) {
	e.readDeadline = time.Now().Add(e.readTimeout)
	e.maxReadDeadline = time.Now().Add(e.maxReadTimeout)

	for n > 0 {
		n--
		event := e.events[n]
		fd := event.Fd

		if event.Events & syscall.EPOLLRDHUP != 0 {
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

		time.Sleep(e.maxReadTimeout)
	}
}

// Start 启动连接管理器
func (e *ConnMgr) Start() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	e.dispatcher.Start()

	e.keepAliveMgr.Start()

	go e.StartConnSignalHookWorkers()

	for {
		n, err := e.Wait()
		if err != nil {
			glog.Error("epoll wait error:", err)
			continue
		}

		e.BatchCommit(n)
	}
}

// Stop 结束连接管理器
func (e *ConnMgr) Stop() {
	n := e.connShards.Count()
	for conn := range e.connShards.ValuesIter(n) {
		conn.Close()
	}

	syscall.Close(e.epfd)
}

// ReadDeadline 头部超时时间
func (e *ConnMgr) ReadDeadline() time.Time {
	return e.readDeadline
}

// MaxReadDeadline 主体超时时间
func (e *ConnMgr) MaxReadDeadline() time.Time {
	return e.maxReadDeadline
}

// StartConnSignalHookWorkers 启动连接信号钩子消费者工作池
func (e *ConnMgr) StartConnSignalHookWorkers() {
	for i := 0; i < len(e.connSignalQueue); i++ {
		for j := 0; j < gconf.Config.WorkersPerConnSignalQueue(); j++ {
			go e.StartConnSignalHookWorker(e.connSignalQueue[i])
		}
	}
}

// StartConnSignalHookWorker 启动连接信号钩子消费者
func (m *ConnMgr) StartConnSignalHookWorker(connSignalQueue <- chan trait.ConnSignal) {
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
func (m *ConnMgr) OnConnStart(fn func(conn trait.Connection)) {
	m.connStartHook = fn
}

// OnConnStop 注册连接断开触发的钩子回调
func (m *ConnMgr) OnConnStop(fn func(conn trait.Connection)) {
	m.connStopHook = fn
}

// OnConnNotActive 注册连接不活跃触发的钩子回调
func (m *ConnMgr) OnConnNotActive(fn func(conn trait.Connection)) {
	m.connNotActiveHook = fn
}

// ChooseConnSignalQueue 选择连接信号处理队列
func (m *ConnMgr) ChooseConnSignalQueue(connID uint64) chan <- trait.ConnSignal {
	return m.connSignalQueue[connID%uint64(len(m.connSignalQueue))]
}

// PushConnSignal 推送连接信号
func (m *ConnMgr) PushConnSignal(signal trait.ConnSignal) {
	m.ChooseConnSignalQueue(signal.ID()) <- signal
}
