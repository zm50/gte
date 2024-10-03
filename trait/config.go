package trait

type ServerConfig interface {
	Load(filePath string) error

	ListenIP() string
	ListenPort() int
	NetworkVersion() string
	MaxConns() int
 	MaxPacketSize() int
	EpollTimeout() int
	EpollEventSize() int
	DispatcherQueues() int
	DispatcherQueueLen() int
	WorkersPerDispatcherQueue() int
	TaskQueues() int
	TaskQueueLen() int
	WorkersPerTaskQueue() int

	WithListenIP(string) ServerConfig
	WithListenPort(int) ServerConfig
	WithNetworkVersion(string) ServerConfig
	WithMaxConns(int) ServerConfig
 	WithMaxPacketSize(int) ServerConfig
	WithEpollTimeout(int) ServerConfig
	WithEpollEventSize(int) ServerConfig
	WithDispatcherQueues(int) ServerConfig
	WithDispatcherQueueLen(int) ServerConfig
	WithWorkersPerDispatcherQueue(int) ServerConfig
	WithTaskQueues(int) ServerConfig
	WithTaskQueueLen(int) ServerConfig
	WithWorkersPerTaskQueue(int) ServerConfig
}
