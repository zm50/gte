package gcore

import "github.com/zm50/gte/trait"

// Router 任务执行流路由器
type Router[T any] struct {
	apis map[uint32]trait.TaskFlow[T]
}

var _ trait.Router[any] = (*Router[any])(nil)

// NewRouter 创建一个新的任务流路由器
func NewRouter[T any]() trait.Router[T] {
	return &Router[T]{
		apis: make(map[uint32]trait.TaskFlow[T]),
	}
}

// Regist 注册任务执行逻辑
func (r *Router[T]) Regist(id uint32, flow ...trait.TaskFunc[T]) {
	if _, ok := r.apis[id]; ok {
		r.apis[id] = r.apis[id].Append(flow...)
	} else {
		r.apis[id] = NewTaskFlow(flow...)
	}
}

// RegistFlow 注册一个任务执行执行流
func (r *Router[T]) RegistFlow(id uint32, flow trait.TaskFlow[T]) {
	r.apis[id] = flow
}

// TaskFlow 根据消息ID获取任务执行流
func (r *Router[T]) TaskFlow(id uint32) trait.TaskFlow[T] {
	return r.apis[id].Fork()
}
