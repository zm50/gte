package global

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig gte框架内部配置
type ServerConfig struct {
	MaxConns int
	MaxPacketSize int
	DispatcherQueues int
	DispatcherQueueLen int
	WorkersPerDispatcherQueue int
	TaskQueues int
	TaskQueueLen int
	WorkersPerTaskQueue int
}

// Config gte框架默认配置
var Config = &ServerConfig{
	MaxConns: 1024,
	MaxPacketSize: 4096,

	DispatcherQueues: 10,
	DispatcherQueueLen: 100,
	WorkersPerDispatcherQueue: 10,

	TaskQueues: 10,
	TaskQueueLen: 100,
	WorkersPerTaskQueue: 10,
}

// LoadConfig 从配置文件中加载配置
func LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// 配置文件打开失败
		fmt.Println("config file open failed")
		return err
	}

	err = yaml.Unmarshal(data, Config)
	if err != nil {
		// 配置文件解析失败
		fmt.Println("config file parse failed")
		return err
	}

	return nil
}
