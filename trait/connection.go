package trait

import (
	"net"
	"os"
)

type Socket interface {
	net.Conn
	File() (*os.File, error)
}

// Connection 连接模块抽象层
type Connection[T any] interface {
	Socket

	ID() uint64
	Send(data []byte) error
	SendMsg(msgID uint32, data []byte) error
	Stop()
	BatchCommit() error
	IsAlive() bool
	IsNotActive() bool
	IsClose() bool
	State() uint32
	SetState(state uint32)

	Property() T
	SetProperty(T)
}
