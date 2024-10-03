package trait

type TaskFunc interface {
	Execute(ctx Context)
}

type TaskFlow interface {
	Extend(...TaskFunc)
	Append(fs ...TaskFunc) TaskFlow
	Execute(idx int, ctx Context)
	Len() int
}
