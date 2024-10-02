package gnet

import "github.com/go75/gte/trait"

// Router 任务流路由器
type Router struct {
	apis map[uint16]trait.TaskFlow
}

var _ trait.Router = (*Router)(nil)

// NewRouter 创建一个新的任务流路由器
func NewRouter() trait.Router {
	return &Router{
		apis: make(map[uint16]trait.TaskFlow),
	}
}

// Regist 注册一个任务流
func (r *Router) Regist(id uint16, flow trait.TaskFlow) {
	r.apis[id] = flow
}

// TaskFlow 根据任务流ID获取任务流
func (r *Router) TaskFlow(id uint16) trait.TaskFlow {
	return r.apis[id]
}
