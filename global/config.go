package global

import (
	"github.com/zm50/gte/common/constant"
	"github.com/zm50/gte/core"
	"github.com/zm50/gte/trait"
)

// conf gte框架全局配置
var conf trait.ServerConfig = newDefaultConfig()

func newDefaultConfig() trait.ServerConfig {
	defaultConfig := &core.ServerConfig{}

	return defaultConfig.
	WithListenIP("0.0.0.0").
	WithListenPort(8080).
	WithNetworkVersion("tcp4").
	WithReadTry(1).
	WithWriteInternal(100).
	WithNetworkMode(constant.TCPNetowrkMode).
	WithMaxConns(1024).
	WithMaxPacketSize(4096).
	WithEpollTimeout(-1).
	WithEpollEventSize(128).
	WithDispatcherQueues(8).
	WithDispatcherQueueLen(128).
	WithWorkersPerDispatcherQueue(2).
	WithTaskQueues(8).
	WithTaskQueueLen(128).
	WithWorkersPerTaskQueue(4).
	WithWebsocketQueueLen(16).
	WithConnSignalQueues(2).
	WithConnSignalQueueLen(4).
	WithWorkersPerConnSignalQueue(2).
	WithConnShardCount(16).
	WithHealthCheckInterval(120000).
	WithMessagePoolSize(10240).
	WithLogFilename("./gte.log").
	WithLogMaxSize(100).
	WithLogMaxBackups(100).
	WithLogMaxAge(30).
	WithLogCompress(false)
}

func Config() trait.ServerConfig {
	return conf
}
