package gcore

import "github.com/zm50/gte/trait"

// ConnSignal 连接状态信号
type ConnSignal[T any] struct {
	trait.Connection[T]
	signal uint8
}

var _ trait.ConnSignal[int] = (*ConnSignal[int])(nil)

// NewConnSignal 创建连接状态信号
func NewConnSignal[T any](conn trait.Connection[T], signal uint8) *ConnSignal[T] {
	return &ConnSignal[T]{conn, signal}
}

// Signal 连接状态信号
func (e *ConnSignal[T]) Signal() uint8 {
	return e.signal
}
