package trait

type Request[T any] interface {
	Message

	Conn() Connection[T]
}
