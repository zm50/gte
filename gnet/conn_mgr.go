package gnet

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/go75/gte/trait"
)

// ConnMgr 连接管理模块，管理客户端的连接，监听待读取数据的连接，并将连接提交给下游的消息分发模块进行处理
type ConnMgr struct {
	fd int

	timeout int

	events []syscall.EpollEvent

	headerTimeout time.Duration
	bodyTimeout time.Duration

	taskMgr trait.TaskMgr

	// key: fd, value: Conn
	conns map[int32]trait.Connection
	connsLock sync.RWMutex
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
		conns: make(map[int32]trait.Connection),
		connsLock: sync.RWMutex{},
	}

	return connMgr, nil
}

// Get 在连接管理器中查群连接
func (e *ConnMgr) Get(fd int32) trait.Connection {
	e.connsLock.RLock()
	defer e.connsLock.RUnlock()

	return e.conns[fd]
}

// Add 在连接管理器中添加连接
func (e *ConnMgr) Add(conn trait.Connection) error {
	sock, err := conn.File()
	if err != nil {
		return err
	}
	fd := sock.Fd()

	if _, ok := e.conns[int32(fd)]; ok {
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

	e.conns[int32(fd)] = conn

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

	delete(e.conns, int32(fd))

	return nil
}

// Wait 等待事件发生
func (e *ConnMgr) Wait() (int, error) {
	n, err := syscall.EpollWait(e.fd, e.events, e.timeout)

	return n, err
}

func (e *ConnMgr) BatchProcess(n int) {
	e.connsLock.RLock()
	defer e.connsLock.RUnlock()

	headerDeadline := time.Now().Add(e.headerTimeout)
	bodyDeadline := time.Now().Add(e.bodyTimeout)

	for n > 0 {
		n--
		fd := int32(e.events[n].Fd)
		conn := e.conns[fd]

		err := e.taskMgr.BatchDispatch(conn, headerDeadline, bodyDeadline)
		if err != nil {
			fmt.Println("batch dispatch error:", err)
		}
	}

	// 基于全局的连接管理模块，获取待处理的连接列表
}

// Start 启动连接管理器
func (e *ConnMgr) Start() {
	for {
		n, err := e.Wait()
		if err != nil {
			fmt.Println("epoll wait error:", err)
			continue
		}

		e.BatchProcess(n)
	}
}

// Close 关闭连接管理器
func (e *ConnMgr) Close() {
	for _, conn := range e.conns {
		conn.Close()
	}

	syscall.Close(e.fd)
}
