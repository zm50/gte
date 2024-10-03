package gcore

import (
	"errors"
	"sync"

	"github.com/go75/gte/gpack"
	"github.com/go75/gte/trait"
)

// Connection 连接模块
type Connection struct {
	// 底层连接的套接字
	trait.Socket

	// 连接的唯一标识
	id int32

	isClosed bool
	//防止连接并发写的锁
	writeLock *sync.Mutex
	exitCh   chan struct{}
}

var _ trait.Connection = (*Connection)(nil)

// NewConnection 创建一个新的连接对象
func NewConnection(connID int32, socket trait.Socket) *Connection {
	return &Connection{
		id: connID,
		Socket: socket,
		isClosed: false,
		exitCh: make(chan struct{}, 1),
	}
}

// ID 返回连接ID
func (c *Connection) ID() int32 {
	return c.id
}

// Send 发送消息给客户端
func (c *Connection) Send(msgID uint16, data []byte) error {
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
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Socket.Close()

	close(c.exitCh)
}
