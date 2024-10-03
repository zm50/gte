package gcore

import (
	"github.com/go75/gte/constant"
	"github.com/go75/gte/trait"
)

// Context 任务上下文
type Context struct {
	trait.Request

	taskIdx int
	tasks trait.TaskFlow
}

var _ trait.Context = (*Context)(nil)

// NewContext 创建任务上下文
func NewContext(req trait.Request, handlers trait.TaskFlow) *Context {
	return &Context{
		Request:  req,

		taskIdx: -1,
		tasks: handlers,
	}
}

// Next 执行下一个任务
func (c *Context) Next() {
	c.taskIdx++
	if c.taskIdx < c.tasks.Len() {
		c.tasks.Execute(c.taskIdx, c)
		c.taskIdx++
	}
}

// Abort 中止任务流
func (c *Context) Abort() {
	c.taskIdx = constant.AbortIndex
}
