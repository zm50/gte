package gcore

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/glog"
	"github.com/zm50/gte/gpack"
	"github.com/zm50/gte/trait"
)

// TCPConnection TCP连接模块
type TCPConnection[T any] struct {
	// 连接的唯一标识
	id uint64

	property T

	// 底层连接的套接字
	trait.Socket

	state *atomic.Uint32
	//防止连接并发写的锁
	writeLock sync.Mutex

	connMgr   trait.ConnMgr[T]
	taskMgr   trait.TaskMgr[T]
	taskQueue chan<- trait.Request[T]
}

var _ trait.Connection[int] = (*TCPConnection[int])(nil)

// NewTCPConnection 创建一个新的连接对象
func NewTCPConnection[T any](connID uint64, socket trait.Socket, connMgr trait.ConnMgr[T], taskMgr trait.TaskMgr[T]) trait.Connection[T] {
	state := &atomic.Uint32{}
	state.Store(constant.ConnActiveState)

	conn := &TCPConnection[T]{
		id:        connID,
		Socket:    socket,
		state:     state,
		writeLock: sync.Mutex{},
		connMgr:   connMgr,
		taskMgr:   taskMgr,
		taskQueue: taskMgr.ChooseQueue(connID),
	}

	return conn
}

// ID 返回连接ID
func (c *TCPConnection[T]) ID() uint64 {
	return c.id
}

// Send 发送数据给客户端
func (c *TCPConnection[T]) Send(data []byte) error {
	if !c.IsActive() {
		return errors.New("connection is closed when send message")
	}

	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	_, err := c.Socket.Write(data)
	if err != nil {
		glog.Error("send data to conn %d err: %v", c.id, err)
		return err
	}

	return nil
}

// SendMsg 发送消息给客户端
func (c *TCPConnection[T]) SendMsg(msgID uint32, data []byte) error {
	//封装message消息
	message := gpack.NewMessage(msgID, data)

	//封包
	response := gpack.PackTCP(message)

	err := c.Send(response)

	return err
}

// Stop 关闭连接
func (c *TCPConnection[T]) Stop() {
	if !c.IsActive() {
		return
	}

	c.SetState(constant.ConnCloseState)
	c.Socket.Close()
}

// BatchCommit 批量提交消息
func (c *TCPConnection[T]) BatchCommit() error {
	for time.Now().Before(c.connMgr.ReadDeadline()) {
		header := make([]byte, 8)

		// 设置header读取超时时间
		c.SetReadDeadline(c.connMgr.ReadDeadline())

		_, err := io.ReadFull(c, header)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 数据包读取超时
				return nil
			}

			return err
		}

		// 设置body读取超时时间
		c.SetReadDeadline(c.connMgr.ReadDeadline())

		msg, err := gpack.UnpackTCP(c)
		if err != nil {
			glog.Error("unpack tcp message err:", err)
			return err
		}

		// 提交消息，处理数据
		request := NewRequest(c, msg)

		c.taskQueue <- request
	}

	return nil
}

// IsActive 连接是否活跃
func (c *TCPConnection[T]) IsActive() bool {
	return c.state.Load() == constant.ConnActiveState
}

// IsNotActive 连接是否不活跃
func (c *TCPConnection[T]) IsNotActive() bool {
	return c.state.Load() == constant.ConnNotActiveState
}

// IsInspect 连接是否处于检查状态
func (c *TCPConnection[T]) IsInspect() bool {
	return c.state.Load() == constant.ConnInspectState
}

// IsClose 连接是否关闭
func (c *TCPConnection[T]) IsClose() bool {
	return c.state.Load() == constant.ConnCloseState
}

// State 获取连接状态
func (c *TCPConnection[T]) State() uint32 {
	return c.state.Load()
}

// SetState 设置连接状态
func (c *TCPConnection[T]) SetState(state uint32) {
	c.state.Store(state)
}

// Property 获取连接属性
func (c *TCPConnection[T]) Property() T {
	return c.property
}

// SetProperty 设置连接属性
func (c *TCPConnection[T]) SetProperty(property T) {
	c.property = property
}

// Websocket websocket连接
type WebsocketConnection[T any] struct {
	// 连接的唯一标识
	id uint64

	property T

	*websocket.Conn

	//防止连接并发写的锁
	writeLock sync.Mutex

	state *atomic.Uint32

	connMgr   trait.ConnMgr[T]
	taskMgr   trait.TaskMgr[T]
	taskQueue chan<- trait.Request[T]
}

var _ trait.Connection[int] = (*WebsocketConnection[int])(nil)

func NewWebsocketConnection[T any](connID uint64, conn *websocket.Conn, connMgr trait.ConnMgr[T], taskMgr trait.TaskMgr[T]) trait.Connection[T] {
	state := &atomic.Uint32{}
	state.Store(constant.ConnActiveState)

	return &WebsocketConnection[T]{
		id:        connID,
		Conn:      conn,
		writeLock: sync.Mutex{},
		state:     state,
		connMgr:   connMgr,
		taskMgr:   taskMgr,
		taskQueue: taskMgr.ChooseQueue(connID),
	}
}

func (w *WebsocketConnection[T]) Read(b []byte) (n int, err error) {
	messageType, data, err := w.Conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	if messageType == websocket.BinaryMessage {
		w.SetState(constant.ConnActiveState)
		return 0, errors.New("not support message type")
	}

	return copy(b, data), nil
}

func (w *WebsocketConnection[T]) Write(b []byte) (int, error) {
	err := w.Conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (w *WebsocketConnection[T]) Close() error {
	return w.Conn.Close()
}

func (w *WebsocketConnection[T]) LocalAddr() net.Addr {
	return w.Conn.LocalAddr()
}

func (w *WebsocketConnection[T]) RemoteAddr() net.Addr {
	return w.Conn.RemoteAddr()
}

func (w *WebsocketConnection[T]) SetDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebsocketConnection[T]) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebsocketConnection[T]) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

func (w *WebsocketConnection[T]) File() (*os.File, error) {
	return nil, errors.New("not support file")
}

func (w *WebsocketConnection[T]) ID() uint64 {
	return w.id
}

// Send 发送数据给客户端
func (w *WebsocketConnection[T]) Send(data []byte) error {
	if !w.IsActive() {
		return errors.New("connection is closed when send message")
	}

	w.writeLock.Lock()
	defer w.writeLock.Unlock()

	_, err := w.Write(data)
	if err != nil {
		glog.Error("send data to conn %d err: %v", w.id, err)
		return err
	}

	return nil
}

// SendMsg 发送消息给客户端
func (w *WebsocketConnection[T]) SendMsg(msgID uint32, data []byte) error {
	//封装message消息
	message := gpack.NewMessage(msgID, data)

	//封包
	response := gpack.PackWebsocket(message)

	err := w.Send(response)

	return err
}

func (w *WebsocketConnection[T]) Stop() {
	if !w.IsActive() {
		return
	}

	w.SetState(constant.ConnCloseState)
	w.Close()
}

func (w *WebsocketConnection[T]) BatchCommit() error {
	w.SetReadDeadline(w.connMgr.MaxReadDeadline())

	for time.Now().Before(w.connMgr.MaxReadDeadline()) {
		messageType, data, err := w.Conn.ReadMessage()
		if err != nil {
			glog.Error("read websocket message err:", err)
			return err
		}

		if messageType == websocket.BinaryMessage {
			w.SetState(constant.ConnActiveState)
			continue
		}

		msg, err := gpack.UnpackWebsocket(data)
		if err != nil {
			glog.Error("unpack websocket message err:", err)
			return err
		}

		request := NewRequest(w, msg)

		w.taskQueue <- request
	}

	return nil
}

// IsActive 连接是否活跃
func (w *WebsocketConnection[T]) IsActive() bool {
	return w.state.Load() == constant.ConnActiveState
}

// IsNotActive 连接是否不活跃
func (w *WebsocketConnection[T]) IsNotActive() bool {
	return w.state.Load() == constant.ConnNotActiveState
}

// IsInspect 连接是否处于检查状态
func (w *WebsocketConnection[T]) IsInspect() bool {
	return w.state.Load() == constant.ConnInspectState
}

// IsClose 连接是否关闭
func (w *WebsocketConnection[T]) IsClose() bool {
	return w.state.Load() == constant.ConnCloseState
}

// State 获取连接状态
func (w *WebsocketConnection[T]) State() uint32 {
	return w.state.Load()
}

// SetState 设置连接状态
func (w *WebsocketConnection[T]) SetState(state uint32) {
	w.state.Store(state)
}

// Property 获取连接属性
func (w *WebsocketConnection[T]) Property() T {
	return w.property
}

// SetProperty 设置连接属性
func (w *WebsocketConnection[T]) SetProperty(property T) {
	w.property = property
}
