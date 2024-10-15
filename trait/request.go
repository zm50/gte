package trait

type Request interface {
	Message

	ConnID() uint64
}
