package trait

type TaskFunc interface {
	Execute(ctx Context)
}

type TaskFlow interface {
	Append(fs ...TaskFunc) TaskFlow
	Fork() TaskFlow
	Execute(idx int, ctx Context)
	Len() int
	Funcs() []TaskFunc
}
