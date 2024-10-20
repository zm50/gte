package gpack

import (
	"github.com/zm50/gte/trait"
)

// Message 业务消息
type Message struct {
	//消息ID
	id uint32
	//消息的长度
	dataLen uint32
	//消息的内容
	data []byte
}

var _ trait.Message = (*Message)(nil)

// 创建一个message
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

// ID 返回消息ID
func (m *Message) ID() uint32 {
	return m.id
}

// DataLen 返回消息体的长度
func (m *Message) DataLen() uint32 {
	return m.dataLen
}

// Data 返回消息的内容
func (m *Message) Data() []byte {
	return m.data
}

// SetID 设置消息ID
func (m *Message) SetID(id uint32) {
	m.id = id
}

// SetDataLen 设置消息体的长度
func (m *Message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.data = data
}
