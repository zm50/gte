package gcore

import "github.com/go75/gte/trait"

// RouterGroup 路由组
type RouterGroup struct {
	trait.Router

	baseTaskFlow trait.TaskFlow
}

var _ trait.RouterGroup = (*RouterGroup)(nil)

// NewRouterGroup 创建路由组
func NewRouterGroup(engine trait.Router) trait.RouterGroup {
	return &RouterGroup{
		Router:   NewRouter(),
		baseTaskFlow: NewTaskFlow(),
	}
}

// Group 子路由组
func (g *RouterGroup) Group(flow ...trait.TaskFunc) trait.RouterGroup {
	group := &RouterGroup{
		Router:   g.Router,
		baseTaskFlow: g.baseTaskFlow.Append(flow...),
	}

	return group
}

// Use 注册插件
func (g *RouterGroup) Use(flow ...trait.TaskFunc) {
	g.baseTaskFlow = g.baseTaskFlow.Append(flow...)
}

// Regist 注册任务执行逻辑
func (g *RouterGroup) Regist(id uint16, flow ...trait.TaskFunc) {
	g.Router.Regist(id, g.baseTaskFlow.Append(flow...).Funcs()...)
}

// RegistFlow 注册任务执行流
func (g *RouterGroup) RegistFlow(id uint16, flow trait.TaskFlow) {
	g.Router.RegistFlow(id, flow)
}
