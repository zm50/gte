package gnet

import "github.com/go75/gte/trait"

// TaskFunc 任务处理函数
type TaskFunc func(trait.Context)

// Execute 执行任务
func (h TaskFunc) Execute(ctx trait.Context) {
	h(ctx)
}

var _ trait.TaskFunc = TaskFunc(nil)

// TaskFlow 任务流
type TaskFlow []trait.TaskFunc

var _ trait.TaskFlow = (*TaskFlow)(nil)

// NewTaskFlow 创建任务流
func NewTaskFlow(fs ...trait.TaskFunc) trait.TaskFlow {
	flow := TaskFlow(fs)
	return &flow
}

// Extend 扩展任务流
func (h *TaskFlow) Extend(fs ...trait.TaskFunc) {
	*h = append(*h, fs...)
}

// Append 追加任务流
func (h *TaskFlow) Append(fs ...trait.TaskFunc) trait.TaskFlow {
	flow := make([]trait.TaskFunc, h.Len() + len(fs))
	copy(flow[:h.Len()], *h)
	copy(flow[h.Len():], fs)

	return NewTaskFlow(flow...)
}

// Execute 执行任务
func (h *TaskFlow) Execute(idx int, ctx trait.Context) {
	(*h)[idx].Execute(ctx)
}

// Len 任务流长度
func (h *TaskFlow) Len() int {
	return len(*h)
}
