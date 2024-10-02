package trait

type Context interface {
	Message

	Next()
	Abort()
}
