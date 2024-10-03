package trait

type Request interface {
	Message

	ConnID() uint64
	Send(msgID uint16, data []byte) error
}
