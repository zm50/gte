package trait

type TaskFunc interface {
	Execute(ctx Context)
}

type TaskFlow interface {
	Extend(...TaskFunc)
	Execute(idx int, ctx Context)
	Len() int
}
