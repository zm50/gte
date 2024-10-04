package gcore

import "github.com/go75/gte/trait"

// ConnSignal 连接状态信号
type ConnSignal struct {
	trait.Connection
	signal uint8
}

var _ trait.ConnSignal = (*ConnSignal)(nil)

// NewConnSignal 创建连接状态信号
func NewConnSignal(conn trait.Connection, signal uint8) *ConnSignal {
	return &ConnSignal{conn, signal}
}

// Signal 连接状态信号
func (e *ConnSignal) Signal() uint8 {
	return e.signal
}
