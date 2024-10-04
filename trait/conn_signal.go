package trait

type ConnSignal interface {
	Connection
	Signal() uint8
}
