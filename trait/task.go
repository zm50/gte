package trait

type TaskFunc[T any] interface {
	Execute(ctx Context[T])
}

type TaskFlow[T any] interface {
	Append(fs ...TaskFunc[T]) TaskFlow[T]
	Fork() TaskFlow[T]
	Execute(idx int, ctx Context[T])
	Len() int
	Funcs() []TaskFunc[T]
}
