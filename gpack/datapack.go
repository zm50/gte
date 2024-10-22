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
	binary.LittleEndian.PutUint32(data[:4], msg.DataLen())

	//2.将message的id写入res中
	binary.LittleEndian.PutUint32(data[4:8], msg.ID())

	//3.将message的内容写到res中
	copy(data[8:], msg.Data()[:msg.DataLen()])

	return data
}

// PackWebsocket 将Message封包成Websocket数据流
func PackWebsocket(msg trait.Message) []byte {
	data := make([]byte, msg.DataLen()+4)

	//1.将message的id写入res中
	binary.LittleEndian.PutUint32(data[:4], msg.ID())

	//2.将message的内容写到res中
	copy(data[4:], msg.Data()[:msg.DataLen()])

	return data
}

// UnpackTCPHeader 从TCP连接中读取消息头部
func UnpackTCPHeader(msg trait.Message, reader io.Reader) error {
	header := make([]byte, 8)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return err
	}

	// read data  len (4 byte) and id (4 bytes)
	msg.SetDataLen(binary.LittleEndian.Uint32(header[:4]))
	msg.SetID(binary.LittleEndian.Uint32(header[4:8]))

	return nil
}

// UnpackTCPBody 从TCP连接中读取消息体
func UnpackTCPBody(msg trait.Message, reader io.Reader) error {
	data := make([]byte, msg.DataLen())
	n, err := io.ReadFull(reader, data)
	if err != nil || n != int(msg.DataLen()) {
		return errors.Wrap(err, "read data error")
	}

	msg.ResetData(data...)

	return nil
}

// UnpackWebsocket 从websocket连接中读取消息
func UnpackWebsocket(msg trait.Message, data []byte) error {
	if len(data) < 4 {
		return errors.New("data too short")
	}

	// id (4 bytes)
	msg.SetID(binary.LittleEndian.Uint32(data[:4]))

	msg.SetDataLen(uint32(len(data) - 4))

	msg.ResetData(data[4:]...)

	return nil
}
