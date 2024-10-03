package gnet

import "github.com/go75/gte/trait"

// RouterGroup 路由组
type RouterGroup struct {
	engine      trait.Engine
	baseTaskFlow trait.TaskFlow
}

var _ trait.RouterGroup = (*RouterGroup)(nil)

// NewRouterGroup 创建路由组
func NewRouterGroup(engine trait.Engine) trait.RouterGroup {
	return &RouterGroup{
		engine:      engine,
		baseTaskFlow: NewTaskFlow(),
	}
}

// Group 子路由组
func (g *RouterGroup) Group(flow ...trait.TaskFunc) trait.RouterGroup {
	group := &RouterGroup{
		engine:      g.engine,
		baseTaskFlow: g.baseTaskFlow.Append(flow...),
	}

	return group
}

// Use 注册中间件
func (g *RouterGroup) Use(flow ...trait.TaskFunc) {
	g.baseTaskFlow = g.baseTaskFlow.Append(flow...)
}

// Regist 注册任务流
func (g *RouterGroup) Regist(id uint16, flow ...trait.TaskFunc) {
	g.RegistFlow(id, g.baseTaskFlow.Append(flow...))
}

// RegistFlow 注册任务流
func (g *RouterGroup) RegistFlow(id uint16, flow trait.TaskFlow) {
	g.engine.RegistFlow(id, flow)
}
