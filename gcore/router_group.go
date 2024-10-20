package gcore

import (
	"github.com/zm50/gte/trait"
)

// RouterGroup 路由组
type RouterGroup[T any] struct {
	trait.Router[T]

	baseTaskFlow trait.TaskFlow[T]
}

var _ trait.RouterGroup[any] = (*RouterGroup[any])(nil)

// NewRouterGroup 创建路由组
func NewRouterGroup[T any](rootRouter trait.Router[T]) trait.RouterGroup[T] {
	return &RouterGroup[T]{
		Router:       rootRouter,
		baseTaskFlow: NewTaskFlow[T](),
	}
}

// Group 子路由组
func (g *RouterGroup[T]) Group(flow ...trait.TaskFunc[T]) trait.RouterGroup[T] {
	group := &RouterGroup[T]{
		Router:       g.Router,
		baseTaskFlow: g.baseTaskFlow.Append(flow...),
	}

	return group
}

// Use 注册插件
func (g *RouterGroup[T]) Use(flow ...trait.TaskFunc[T]) {
	g.baseTaskFlow = g.baseTaskFlow.Append(flow...)
}

// Regist 注册任务执行逻辑
func (g *RouterGroup[T]) Regist(id uint32, flow ...trait.TaskFunc[T]) {
	g.Router.Regist(id, g.baseTaskFlow.Append(flow...).Funcs()...)
}

// RegistFlow 注册任务执行流
func (g *RouterGroup[T]) RegistFlow(id uint32, flow trait.TaskFlow[T]) {
	g.Router.RegistFlow(id, flow)
}
