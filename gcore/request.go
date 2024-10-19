package gcore

import "github.com/go75/gte/trait"

// Request 请求对象
type Request struct {
	trait.Connection
	trait.Message
}

var _ trait.Request = (*Request)(nil)

// NewRequest 创建请求对象
func NewRequest(conn trait.Connection, msg trait.Message) trait.Request {
	return &Request{
		Connection: conn,
		Message:    msg,
	}
}

// ConnID 连接ID
func (r *Request) ConnID() uint64 {
	return r.Connection.ID()
}

// ID 消息ID
func (r *Request) ID() uint32 {
	return r.Message.ID()
}
