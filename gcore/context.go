package gcore

import (
	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/trait"
)

// Context 任务上下文
type Context[T any] struct {
	trait.Request[T]

	taskIdx int
	tasks trait.TaskFlow[T]
}

var _ trait.Context[int] = (*Context[int])(nil)

// NewContext 创建任务上下文
func NewContext[T any](req trait.Request[T], handlers trait.TaskFlow[T]) *Context[T] {
	return &Context[T]{
		Request:  req,

		taskIdx: -1,
		tasks: handlers,
	}
}

// Next 执行下一个任务
func (c *Context[T]) Next() {
	c.taskIdx++
	if c.taskIdx < c.tasks.Len() {
		c.tasks.Execute(c.taskIdx, c)
		c.taskIdx++
	}
}

// Abort 中止任务流
func (c *Context[T]) Abort() {
	c.taskIdx = constant.AbortIndex
}
