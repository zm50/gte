package trait

type ServerConfig interface {
	Load(filePath string) error
	Export(filePath string) error

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
	WorkersPerDispatcherQueue() int
	TaskQueues() int
	TaskQueueLen() int
	WorkersPerTaskQueue() int
	WebsocketQueueLen() int
	ConnSignalQueues() int
	ConnSignalQueueLen() int
	WorkersPerConnSignalQueue() int
	ConnShardCount() int
	HealthCheckInterval() int

	WithListenIP(string) ServerConfig
	WithListenPort(int) ServerConfig
	WithNetworkVersion(string) ServerConfig
	WithReadTimeout(int) ServerConfig
	WithMaxReadTimeout(int) ServerConfig
	WithNetworkMode(int) ServerConfig
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
	WithWebsocketQueueLen(int) ServerConfig
	WithConnSignalQueues(int) ServerConfig
	WithConnSignalQueueLen(int) ServerConfig
	WithWorkersPerConnSignalQueue(int) ServerConfig
	WithConnShardCount(connShardCount int) ServerConfig
	WithHealthCheckInterval(int) ServerConfig
}
