package trait

type Message interface {
	ID() uint16
	DataLen() uint16
	Data() []byte

	SetID(uint16)
	SetDataLen(uint16)
	SetData([]byte)
}
