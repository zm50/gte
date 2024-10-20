package gcore

import (
	"github.com/zm50/gte/gconf"
	"github.com/zm50/gte/trait"
)

// TaskMgr 任务管理器
type TaskMgr[T any] struct {
	trait.RouterGroup[T]

	taskQueues []chan trait.Request[T]
}

var _ trait.TaskMgr[any] = (*TaskMgr[any])(nil)

// NewTaskMgr 创建任务管理器
func NewTaskMgr[T any]() trait.TaskMgr[T] {
	taskQueues := make([]chan trait.Request[T], gconf.Config.TaskQueues())
	for i := 0; i < len(taskQueues); i++ {
		taskQueues[i] = make(chan trait.Request[T], gconf.Config.TaskQueueLen())
	}

	// 新建任务处理路由器与分组路由
	rootRouter := NewRouter[T]()
	routerGroup := NewRouterGroup(rootRouter)

	return &TaskMgr[T]{
		RouterGroup: routerGroup,
		taskQueues:  taskQueues,
	}
}

// Start 启动任务管理器
func (m *TaskMgr[T]) Start() {
	for i := 0; i < len(m.taskQueues); i++ {
		for j := 0; j < gconf.Config.WorkersPerTaskQueue(); j++ {
			go m.StartWorker(m.taskQueues[i])
		}
	}
}

// StartWorker 启动任务消费者
func (m *TaskMgr[T]) StartWorker(taskQueue <-chan trait.Request[T]) {
	for request := range taskQueue {
		flow := m.TaskFlow(request.ID())
		ctx := NewContext(request, flow)
		ctx.Next()
	}
}

// ChooseQueue 选择处理连接的队列
func (m *TaskMgr[T]) ChooseQueue(connID uint64) chan<- trait.Request[T] {
	// 负载均衡，选择队列
	return m.taskQueues[connID%uint64(len(m.taskQueues))]
}

// Submit 提交任务
func (m *TaskMgr[T]) Submit(request trait.Request[T]) {
	m.ChooseQueue(request.Conn().ID()) <- request
}

// Use 注册插件
func (m *TaskMgr[T]) Use(flow ...trait.TaskFunc[T]) {
	m.RouterGroup.Use(flow...)
}

// Regist 注册任务流
func (m *TaskMgr[T]) Regist(id uint32, flow ...trait.TaskFunc[T]) {
	m.RouterGroup.Regist(id, flow...)
}

// Regist 注册任务流
func (m *TaskMgr[T]) RegistFlow(id uint32, flow trait.TaskFlow[T]) {
	m.RouterGroup.RegistFlow(id, flow)
}
