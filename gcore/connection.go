package gcore

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go75/gte/gpack"
	"github.com/go75/gte/trait"
	"github.com/gorilla/websocket"
)

// TCPConnection TCP连接模块
type TCPConnection struct {
	// 连接的唯一标识
	id uint64

	// 底层连接的套接字
	trait.Socket

	isClosed bool
	//防止连接并发写的锁
	writeLock sync.Mutex
	exitCh   chan struct{}

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
	taskQueue chan <- trait.Request
}

var _ trait.Connection = (*TCPConnection)(nil)

// NewTCPConnection 创建一个新的连接对象
func NewTCPConnection(connID uint64, socket trait.Socket, connMgr trait.ConnMgr, taskMgr trait.TaskMgr) trait.Connection {
	conn := &TCPConnection{
		id: connID,
		Socket: socket,
		isClosed: false,
		writeLock: sync.Mutex{},
		exitCh: make(chan struct{}, 1),
		connMgr: connMgr,
		taskMgr: taskMgr,
		taskQueue: taskMgr.ChooseQueue(connID),
	}

	return conn
}

// ID 返回连接ID
func (c *TCPConnection) ID() uint64 {
	return c.id
}

// Send 发送消息给客户端
func (c *TCPConnection) Send(msgID uint16, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed when send message")
	}

	//封装message消息
	message := gpack.NewMessage(msgID, data)

	//封包
	response := gpack.Pack(message)

	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	_, err := c.Socket.Write(response)

	return err
}

// Stop 关闭连接
func (c *TCPConnection) Stop() {
	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Socket.Close()

	close(c.exitCh)
}

// BatchCommit 批量提交消息
func (c *TCPConnection) BatchCommit() error {
	for time.Now().Before(c.connMgr.ReadDeadline()) {
		header := make([]byte, 4)

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

		id, dataLen := gpack.UnpackHeader(header)
		body := make([]byte, dataLen)
		_, err = io.ReadFull(c, body)
		if err != nil {
			return err
		}

		msg := gpack.NewMessage(id, body)

		// 提交消息，处理数据		
		request := NewRequest(c, msg)

		c.taskQueue <- request
	}

	return nil
}

// Websocket websocket连接
type WebsocketConnection struct {
	// 连接的唯一标识
	id uint64

	*websocket.Conn

	//防止连接并发写的锁
	writeLock sync.Mutex

	isActive atomic.Bool

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
	taskQueue chan <- trait.Request
}

var _ trait.Socket = (*WebsocketConnection)(nil)

func NewWebsocketConnection(connID uint64, conn *websocket.Conn, connMgr trait.ConnMgr, taskMgr trait.TaskMgr) trait.Connection {
	isActive := atomic.Bool{}
	isActive.Store(true)

	return &WebsocketConnection{
		id: connID,
		Conn: conn,
		writeLock: sync.Mutex{},
		isActive: atomic.Bool{},
		connMgr: connMgr,
		taskMgr: taskMgr,
		taskQueue: taskMgr.ChooseQueue(connID),
	}
}

func (w *WebsocketConnection) Read(b []byte) (n int, err error) {
	messageType, data, err := w.Conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	if messageType == websocket.BinaryMessage {
		w.isActive.Store(true)
		return 0, errors.New("not support message type")
	}

	return copy(b, data), nil
}

func (w *WebsocketConnection) Write(b []byte) (int, error) {
	err := w.Conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (w *WebsocketConnection) Close() error {
	return w.Conn.Close()
}

func (w *WebsocketConnection) LocalAddr() net.Addr {
	return w.Conn.LocalAddr()
}

func (w *WebsocketConnection) RemoteAddr() net.Addr {
	return w.Conn.RemoteAddr()
}

func (w *WebsocketConnection) SetDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebsocketConnection) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebsocketConnection) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

func (w *WebsocketConnection) File() (*os.File, error) {
	return nil, errors.New("not support file")
}

func (w *WebsocketConnection) ID() uint64 {	
	return w.id
}

func (w *WebsocketConnection) Send(msgID uint16, data []byte) error {
	if !w.isActive.Load() {
		return errors.New("connection is closed when send message")
	}

	//封装message消息
	message := gpack.NewMessage(msgID, data)

	//封包
	response := gpack.Pack(message)

	w.writeLock.Lock()
	defer w.writeLock.Unlock()
	err := w.Conn.WriteMessage(websocket.BinaryMessage, response)

	return err
}

func (w *WebsocketConnection) Stop() {
	w.isActive.Store(false)
	w.Conn.Close()
}

func (w *WebsocketConnection) BatchCommit() error {
	w.SetReadDeadline(w.connMgr.MaxReadDeadline())

	for time.Now().Before(w.connMgr.MaxReadDeadline()) {
		messageType, data, err := w.Conn.ReadMessage()
		if err != nil {
			return err
		}
	
		if messageType == websocket.BinaryMessage {
			w.isActive.Store(true)
			continue
		}

		msg := gpack.UnpackFullData(data)

		request := NewRequest(w, msg)

		w.taskQueue <- request
	}

	return nil
}
