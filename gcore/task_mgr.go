package gcore

import (
	"github.com/go75/gte/gconf"
	"github.com/go75/gte/trait"
)

// TaskMgr 任务管理器
type TaskMgr struct {
	trait.RouterGroup

	taskQueues []chan trait.Request
}

var _ trait.TaskMgr = (*TaskMgr)(nil)

// NewTaskMgr 创建任务管理器
func NewTaskMgr() trait.TaskMgr {
	taskQueues := make([]chan trait.Request, gconf.Config.TaskQueues())
	for i := 0; i < len(taskQueues); i++ {
		taskQueues[i] = make(chan trait.Request, gconf.Config.TaskQueueLen())
	}

	// 新建任务处理路由器与分组路由
	router := NewRouter()
	routerGroup := NewRouterGroup(router)

	return &TaskMgr{
		RouterGroup: routerGroup,
		taskQueues: taskQueues,
	}
}

// Start 启动任务管理器
func (m *TaskMgr) Start() {
	for i := 0; i < len(m.taskQueues); i++ {
		for j := 0; j < gconf.Config.WorkersPerTaskQueue(); j++ {
			go m.StartWorker(m.taskQueues[i])
		}
	}
}

// StartWorker 启动任务消费者
func (m *TaskMgr) StartWorker(taskQueue chan trait.Request) {
	for request := range taskQueue {
		flow := m.TaskFlow(request.ID())
		ctx := NewContext(request, flow)
		ctx.Next()
	}
}

// ChooseQueue 选择处理连接的队列
func (m *TaskMgr) ChooseQueue(connID uint64) chan <- trait.Request {
	// 负载均衡，选择队列
	return m.taskQueues[connID % uint64(len(m.taskQueues))]
}

// Submit 提交任务
func (m *TaskMgr) Submit(request trait.Request) {
	m.ChooseQueue(request.ConnID()) <- request
}

// Use 注册插件
func (m *TaskMgr) Use(flow ...trait.TaskFunc) {
	m.RouterGroup.Use(flow...)
}

// Regist 注册任务流
func (m *TaskMgr) Regist(id uint16, flow ...trait.TaskFunc) {
	m.RouterGroup.Regist(id, flow...)
}

// Regist 注册任务流
func (m *TaskMgr) RegistFlow(id uint16, flow trait.TaskFlow) {
	m.RouterGroup.RegistFlow(id, flow)
}
