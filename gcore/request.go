package gcore

import "github.com/zm50/gte/trait"

// Request 请求对象
type Request[T any] struct {
	trait.Connection[T]
	trait.Message
}

var _ trait.Request[any] = (*Request[any])(nil)

// NewRequest 创建请求对象
func NewRequest[T any](conn trait.Connection[T], msg trait.Message) trait.Request[T] {
	return &Request[T]{
		Connection: conn,
		Message:    msg,
	}
}

// Conn 连接
func (r *Request[T]) Conn() trait.Connection[T] {
	return r.Connection
}

// ID 消息ID
func (r *Request[T]) ID() uint32 {
	return r.Message.ID()
}
