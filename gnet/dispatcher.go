package gnet

import (
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/go75/gte/global"
	"github.com/go75/gte/trait"
)

// Dispatcher 消息分发模块，负责读取客户端连接的数据，并对数据进行拆包转换成消息格式，然后分发给下游的任务处理模块对消息进行业务处理
type Dispatcher struct {
	headerDeadline time.Time
	bodyDeadline time.Time

	connQueue []chan trait.Connection
	taskMgr trait.TaskMgr
}

var _ trait.Dispatcher = (*Dispatcher)(nil)

// NewDispatcher 创建一个消息分发器
func NewDispatcher(taskMgr trait.TaskMgr) *Dispatcher {
	connQueue := make([]chan trait.Connection, global.Config.DispatcherQueues)
	for i := 0; i < len(connQueue); i++ {
		connQueue[i] = make(chan trait.Connection, global.Config.DispatcherQueueLen)
	}

	return &Dispatcher{
		connQueue: connQueue,
		taskMgr: taskMgr,
	}
}

// Start 启动消息分发模块
func (d *Dispatcher) Start() {
	for i := 0; i < len(d.connQueue); i++ {
		for j := 0; j < global.Config.DispatcherQueueLen; j++ {
			go d.Dispatch(d.connQueue[i])
		}
	}
}

// StartDispatcher 分发连接数据
func (d *Dispatcher) Dispatch(connQueue chan trait.Connection) {
	// 从conn中读取数据，并将数据提交给taskMgr处理
	for conn := range connQueue {
		d.BatchDispatch(conn)
	}
}

// BatchDispatch 批量分发连接中的数据
func (d *Dispatcher) BatchDispatch(conn trait.Connection) error {
	for time.Now().After(d.headerDeadline) {
		header := make([]byte, 4)

		// 设置header读取超时时间
		conn.SetReadDeadline(d.headerDeadline)

		_, err := io.ReadFull(conn, header)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 数据包读取超时
				return nil
			}

			return err
		}

		// 设置body读取超时时间
		conn.SetReadDeadline(d.bodyDeadline)

		// 读取长度
		dataLen := binary.BigEndian.Uint16(header[2:4])
		// 读取数据
		body := make([]byte, dataLen)	
		_, err = io.ReadFull(conn, body)
		if err != nil {
			return err
		}

		msg := Unpack(header, body)
		// 提交消息，处理数据
		
		request := NewRequest(conn, msg)

		d.taskMgr.Submit(request)
	}

	return nil
}

// SetHeaderDeadline 设置header读取超时时间
func (d *Dispatcher) SetHeaderDeadline(deadline time.Time) {
	d.headerDeadline = deadline
}

// SetBodyDeadline 设置body读取超时时间
func (d *Dispatcher) SetBodyDeadline(deadline time.Time) {
	d.bodyDeadline = deadline
}

// ChooseQueue 选择处理连接的队列
func (d *Dispatcher) ChooseQueue(conn trait.Connection) chan <- trait.Connection {
	// 负载均衡，选择队列
	return d.connQueue[conn.ID() % int32(len(d.connQueue))]
}

// Commit 提交连接到队列
func (d *Dispatcher) Commit(conn trait.Connection) {
	d.ChooseQueue(conn) <- conn
}
