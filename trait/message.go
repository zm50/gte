package trait

type Message interface {
	ID() uint32
	DataLen() uint32
	Data() []byte

	SetID(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
