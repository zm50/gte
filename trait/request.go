package trait

type Request interface {
	Message

	ConnID() int32
	Send(msgID uint16, data []byte) error
}
