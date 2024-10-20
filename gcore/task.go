package gcore

import "github.com/zm50/gte/trait"

// TaskFunc 任务处理函数
type TaskFunc[T any] func(trait.Context[T])

var _ trait.TaskFunc[any] = TaskFunc[any](nil)

// Execute 执行任务
func (h TaskFunc[T]) Execute(ctx trait.Context[T]) {
	h(ctx)
}

// TaskFlow 任务执行流
type TaskFlow[T any] []trait.TaskFunc[T]

var _ trait.TaskFlow[any] = (*TaskFlow[any])(nil)

// NewTaskFlow 创建任务执行流
func NewTaskFlow[T any](fs ...trait.TaskFunc[T]) trait.TaskFlow[T] {
	flow := TaskFlow[T](fs)
	return &flow
}

// Extend 扩展任务执行流
func (h *TaskFlow[T]) Extend(fs ...trait.TaskFunc[T]) {
	*h = append(*h, fs...)
}

// Append 追加任务执行流
func (h *TaskFlow[T]) Append(fs ...trait.TaskFunc[T]) trait.TaskFlow[T] {
	flow := make([]trait.TaskFunc[T], h.Len()+len(fs))
	copy(flow[:h.Len()], *h)
	copy(flow[h.Len():], fs)

	return NewTaskFlow(flow...)
}

// Fork 克隆任务执行流
func (h *TaskFlow[T]) Fork() trait.TaskFlow[T] {
	return h
}

// Execute 执行任务
func (h *TaskFlow[T]) Execute(idx int, ctx trait.Context[T]) {
	(*h)[idx].Execute(ctx)
}

// Len 任务任务执行流
func (h *TaskFlow[T]) Len() int {
	return len(*h)
}

// Funcs 获取所有任务执行逻辑
func (h *TaskFlow[T]) Funcs() []trait.TaskFunc[T] {
	return *h
}

// StatefulFunc 有状态函数
type StatefulFunc[T any] struct {
	Data T
	Func func(trait.Context[T], T)
}

// NewStatefulFunc 创建有状态函数
func NewStatefulFunc[T any](fn func(trait.Context[T], T)) *StatefulFunc[T] {
	return &StatefulFunc[T]{
		Func: fn,
	}
}

// Execute 执行任务
func (f *StatefulFunc[T]) Execute(ctx trait.Context[T]) {
	f.Func(ctx, f.Data)
}

// StatefulFuncFlow 有状态函数流
type StatefulFuncFlow[T any] struct {
	dataProvide func() T
	Data        T
	FuncFlow    []*StatefulFunc[T]
}

// NewStatefulFuncFlow 创建有状态函数流
func NewStatefulFuncFlow[T any](dataProvide func() T) *StatefulFuncFlow[T] {
	return &StatefulFuncFlow[T]{}
}

// Regist 注册有状态函数
func (f *StatefulFuncFlow[T]) Regist(fns ...func(trait.Context[T], T)) {
	for i := 0; i < len(fns); i++ {
		statefulFunc := NewStatefulFunc(fns[i])
		f.FuncFlow = append(f.FuncFlow, statefulFunc)
	}
}

// Append 添加函数到流中
func (f *StatefulFuncFlow[T]) Append(fs ...trait.TaskFunc[T]) trait.TaskFlow[T] {
	flow := make([]trait.TaskFunc[T], f.Len()+len(fs))

	taskFuncs := make([]trait.TaskFunc[T], len(f.FuncFlow))

	for i, fn := range f.FuncFlow {
		taskFuncs[i] = fn
	}

	copy(flow[:f.Len()], taskFuncs)
	copy(flow[f.Len():], fs)

	return f
}

// Fork 复制一个新的有状态函数流
func (f *StatefulFuncFlow[T]) Fork() trait.TaskFlow[T] {
	flow := &StatefulFuncFlow[T]{
		Data: f.dataProvide(),
	}

	statefulFuncs := make([]*StatefulFunc[T], len(f.FuncFlow))

	for i, fn := range f.FuncFlow {
		statefulFuncs[i] = NewStatefulFunc(fn.Func)
	}

	return flow
}

// Execute 执行任务
func (f *StatefulFuncFlow[T]) Execute(idx int, ctx trait.Context[T]) {
	statefulFunc := f.FuncFlow[idx]
	statefulFunc.Data = f.Data
	statefulFunc.Execute(ctx)
}

// Len 获取函数流长度
func (f *StatefulFuncFlow[T]) Len() int {
	return len(f.FuncFlow)
}

// Funcs 获取函数清单
func (f *StatefulFuncFlow[T]) Funcs() []trait.TaskFunc[T] {
	taskFuncs := make([]trait.TaskFunc[T], len(f.FuncFlow))
	for i, fn := range f.FuncFlow {
		taskFuncs[i] = fn
	}

	return taskFuncs
}
