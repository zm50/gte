package gcore

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/go75/gte/core"
	"github.com/go75/gte/trait"
)

// ConnMgr 连接管理模块，管理客户端的连接，监听待读取数据的连接，并将连接提交给下游的消息分发模块进行处理
type ConnMgr struct {
	fd int

	timeout int

	events []syscall.EpollEvent

	headerTimeout time.Duration
	bodyTimeout time.Duration

	dispatcher trait.Dispatcher

	// key: fd, value: Conn
	connShards *core.KVShards[int32, trait.Connection]
}

var _ trait.ConnMgr = (*ConnMgr)(nil)

// NewConnMgr 新建一个连接管理的实例
func NewConnMgr(timeout int, eventSize int) (*ConnMgr, error) {
	// 创建一个epoll句柄
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	// 创建一个连接管理器
	connMgr := &ConnMgr{
		fd: epfd,
		timeout: timeout,
		events: make([]syscall.EpollEvent, eventSize),
		connShards: core.NewKVShards[int32, trait.Connection](32),
	}

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
		return err
	}
	fd := sock.Fd()

	if _, ok := e.Get(int32(fd)); ok {
		return errors.New("connection already exists")
	}

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}
	err = syscall.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, int(fd), &event)
	if err != nil {
		return err
	}

	e.connShards.Set(int32(fd), conn)

	return nil
}

// Del 在连接管理器中删除连接
func (e *ConnMgr) Del(fd int) error {
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}
	err := syscall.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, &event)
	if err != nil {
		return err
	}

	e.connShards.Del(int32(fd))

	return nil
}

// Wait 等待事件发生
func (e *ConnMgr) Wait() (int, error) {
	n, err := syscall.EpollWait(e.fd, e.events, e.timeout)

	return n, err
}

// BatchCommit 批量处理待读取的连接
func (e *ConnMgr) BatchCommit(n int) {
	e.dispatcher.SetHeaderDeadline(time.Now().Add(e.headerTimeout))
	e.dispatcher.SetBodyDeadline(time.Now().Add(e.bodyTimeout))

	for n > 0 {
		n--
		fd := int32(e.events[n].Fd)
		conn, ok := e.Get(fd)
		if !ok {
			fmt.Println("connection not found:", fd)
			continue
		}

		// 提交待读取的连接
		e.dispatcher.Commit(conn)
	}
}

// Start 启动连接管理器
func (e *ConnMgr) Start() {
	for {
		n, err := e.Wait()
		if err != nil {
			fmt.Println("epoll wait error:", err)
			continue
		}

		e.BatchCommit(n)
	}
}

// Close 关闭连接管理器
func (e *ConnMgr) Close() {
	n := e.connShards.Count()
	for conn := range e.connShards.ValuesIter(n) {
		conn.Close()
	}

	syscall.Close(e.fd)
}
