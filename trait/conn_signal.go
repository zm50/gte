package trait

type ConnSignal[T any] interface {
	Connection[T]
	Signal() uint8
}
