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
type Connection interface {
	Socket

	ID() uint64
	Send(msgID uint16, data []byte) error
	Stop()
	BatchCommit() error
}
