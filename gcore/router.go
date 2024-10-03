package gcore

import "github.com/go75/gte/trait"

// Router 任务执行流路由器
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

// Regist 注册任务执行逻辑
func (r *Router) Regist(id uint16, flow ...trait.TaskFunc) {
	if _, ok := r.apis[id]; ok {
		r.apis[id].Extend(flow...)
	} else {
		r.apis[id] = NewTaskFlow(flow...)
	}
}

// RegistFlow 注册一个任务执行执行流
func (r *Router) RegistFlow(id uint16, flow trait.TaskFlow) {
	r.apis[id] = flow
}

// TaskFlow 根据消息ID获取任务执行流
func (r *Router) TaskFlow(id uint16) trait.TaskFlow {
	return r.apis[id]
}
