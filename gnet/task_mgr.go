package gnet

import (
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
