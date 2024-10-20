package trait

type Context[T any] interface {
	Request[T]

	Next()
	Abort()
}
