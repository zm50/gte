package gnet

import (
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/go75/gte/global"
	"github.com/go75/gte/trait"
)

// TaskMgr 任务管理器
type TaskMgr struct {
	taskQueues []chan trait.Request
	router trait.Router
}

var _ trait.TaskMgr = (*TaskMgr)(nil)

// NewTaskMgr 创建任务管理器
func NewTaskMgr(router trait.Router) trait.TaskMgr {
	taskQueues := make([]chan trait.Request, global.Config.TaskQueues)
	for i := 0; i < len(taskQueues); i++ {
		taskQueues[i] = make(chan trait.Request, global.Config.TaskQueueLen)
	}

	return &TaskMgr{
		taskQueues: taskQueues,
		router: router,
	}
}

// Start 启动任务管理器
func (m *TaskMgr) Start() {
	for i := 0; i < len(m.taskQueues); i++ {
		for j := 0; j < global.Config.WorkersPerTaskQueue; j++ {
			go m.StartWorker(m.taskQueues[i])
		}
	}
}

// StartWorker 启动任务消费者
func (m *TaskMgr) StartWorker(taskQueue chan trait.Request) {
	for request := range taskQueue {
		flow := m.router.TaskFlow(request.ID())
		ctx := NewContext(request, flow)
		ctx.Next()
	}
}

// Submit 提交任务
func (m *TaskMgr) Submit(request trait.Request) {
	m.taskQueues[int(request.ConnID()) % len(m.taskQueues)] <- request
}

// Regist 注册任务流
func (m *TaskMgr) Regist(id uint16, flow trait.TaskFlow) {
	m.router.Regist(id, flow)
}

// BatchDispatch 批量分发数据包
func (m *TaskMgr) BatchDispatch(conn trait.Connection, headerDeadline, bodyDeadline time.Time) error {
	for time.Now().After(headerDeadline) {
		header := make([]byte, 4)

		// 设置header读取超时时间
		conn.SetReadDeadline(headerDeadline)

		_, err := io.ReadFull(conn, header)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 数据包读取超时
				return nil
			}

			return err
		}

		// 设置body读取超时时间
		conn.SetReadDeadline(bodyDeadline)

		// 读取长度
		dataLen := binary.BigEndian.Uint16(header[2:4])
		// 读取数据
		body := make([]byte, dataLen)	
		_, err = io.ReadFull(conn, body)
		if err != nil {
			return err
		}

		msg := unpack(header, body)
		// 提交消息，处理数据
		
		request := NewRequest(conn, msg)

		m.Submit(request)
	}

	return nil
}
