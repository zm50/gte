package gpack

import (
	"encoding/binary"

	"github.com/go75/gte/trait"
)

/**
封包，拆包 模块
直接面向TCP连接中的数据流，用于处理TCP粘包
*/

// data开头的2字节是数据的长度,接下来的2字节是数据的id,在接下来是数据的具体内容

// 封包
func Pack(msg trait.Message) []byte {
	data := make([]byte, msg.DataLen()+4)
	
	//1.将message的id写入res中
	binary.BigEndian.PutUint16(data[0:2], msg.ID())
	
	//2.将datalen写到res中
	binary.BigEndian.PutUint16(data[2:4], msg.DataLen())

	//3.将message的内容写到res中
	copy(data[4:], msg.Data()[:msg.DataLen()])

	return data
}

// 拆包
func Unpack(header []byte, body []byte) trait.Message {
	// read data  id (2 byte) and len (2 bytes)
	id := binary.BigEndian.Uint16(header[:2])
	msg := NewMessage(id, body)

	return msg
}
