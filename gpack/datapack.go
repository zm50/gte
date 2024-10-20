package gpack

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
	"github.com/zm50/gte/trait"
)

/**
封包，拆包 模块
直接面向TCP连接中的数据流，用于处理TCP粘包
*/

// data开头的4字节是数据的长度,接下来的4字节是数据的id,在接下来是数据的具体内容

// PackTCP 将Message封包成TCP数据流
func PackTCP(msg trait.Message) []byte {
	data := make([]byte, msg.DataLen()+8)

	//1.将datalen写到res中
	binary.BigEndian.PutUint32(data[:4], msg.DataLen())

	//2.将message的id写入res中
	binary.BigEndian.PutUint32(data[4:8], msg.ID())

	//3.将message的内容写到res中
	copy(data[8:], msg.Data()[:msg.DataLen()])

	return data
}

// PackWebsocket 将Message封包成Websocket数据流
func PackWebsocket(msg trait.Message) []byte {
	data := make([]byte, msg.DataLen()+4)

	//1.将message的id写入res中
	binary.BigEndian.PutUint32(data[:4], msg.ID())

	//2.将message的内容写到res中
	copy(data[4:], msg.Data()[:msg.DataLen()])

	return data
}

// UnpackTCP 从TCP连接中读取数据，解包成Message
func UnpackTCP(reader io.Reader) (trait.Message, error) {
	header := make([]byte, 8)
	n, err := io.ReadFull(reader, header)
	if err != nil || n != 8 {
		return nil, errors.Wrap(err, "read header error")
	}

	// read data  len (4 byte) and id (4 bytes)
	dataLen := binary.BigEndian.Uint32(header[:4])
	id := binary.BigEndian.Uint32(header[4:8])

	data := make([]byte, dataLen)
	n, err = io.ReadFull(reader, data)
	if err != nil || n != int(dataLen) {
		return nil, errors.Wrap(err, "read data error")
	}

	msg := NewMessage(id, data)

	return msg, nil
}

// UnpackWebsocket 基于websocket读取的数据，解包成Message
func UnpackWebsocket(data []byte) (trait.Message, error) {
	if len(data) < 4 {
		return nil, errors.New("data too short")
	}

	// id (4 bytes)
	id := binary.BigEndian.Uint32(data[4:8])

	msg := NewMessage(id, data[4:])

	return msg, nil
}
