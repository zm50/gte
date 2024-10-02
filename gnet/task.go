package gnet

import "github.com/go75/gte/trait"

type TaskFunc func(trait.Context)

func (h TaskFunc) Execute(ctx trait.Context) {
	h(ctx)
}

var _ trait.TaskFunc = TaskFunc(nil)

type TaskFlow []trait.TaskFunc

var _ trait.TaskFlow = (*TaskFlow)(nil)

func NewTaskFlow(fs ...trait.TaskFunc) trait.TaskFlow {
	flow := TaskFlow(fs)
	return &flow
}

func (h *TaskFlow) Extend(fs ...trait.TaskFunc) {
	*h = append(*h, fs...)
}

func (h *TaskFlow) Execute(idx int, ctx trait.Context) {
	(*h)[idx].Execute(ctx)
}

func (h *TaskFlow) Len() int {
	return len(*h)
}
