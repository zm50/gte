package trait

type ServerConfig interface {
	Load(filePath string) error

	ListenIP() string
	ListenPort() int
	NetworkVersion() string
	ReadTimeout() int
	MaxReadTimeout() int
	NetworkMode() int
	MaxConns() int
 	MaxPacketSize() int
	EpollTimeout() int
	EpollEventSize() int
	DispatcherQueues() int
	DispatcherQueueLen() int
	TaskQueues() int
	TaskQueueLen() int
	WorkersPerTaskQueue() int
	WebsocketQueueLen() int

	WithListenIP(string) ServerConfig
	WithListenPort(int) ServerConfig
	WithNetworkVersion(string) ServerConfig
	WithReadTimeout(readTimeout int) ServerConfig
	WithMaxReadTimeout(maxReadTimeout int) ServerConfig
	WithNetworkMode(int) ServerConfig
	WithMaxConns(int) ServerConfig
 	WithMaxPacketSize(int) ServerConfig
	WithEpollTimeout(int) ServerConfig
	WithEpollEventSize(int) ServerConfig
	WithDispatcherQueues(int) ServerConfig
	WithDispatcherQueueLen(int) ServerConfig
	WithTaskQueues(int) ServerConfig
	WithTaskQueueLen(int) ServerConfig
	WithWorkersPerTaskQueue(int) ServerConfig
	WithWebsocketQueueLen(int) ServerConfig
}
