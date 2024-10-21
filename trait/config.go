package trait

type ServerConfig interface {
	Load(filePath string) error
	Export(filePath string) error

	ListenIP() string
	ListenPort() int
	NetworkVersion() string
	ReadTry() int
	WriteInternal() int
	NetworkMode() int
	MaxConns() int32
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
	LogFilename() string
	LogMaxSize() int
	LogMaxBackups() int
	LogMaxAge() int
	LogCompress() bool

	WithListenIP(string) ServerConfig
	WithListenPort(int) ServerConfig
	WithNetworkVersion(string) ServerConfig
	WithReadTry(int) ServerConfig
	WithWriteInternal(int) ServerConfig
	WithNetworkMode(int) ServerConfig
	WithMaxConns(int32) ServerConfig
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
	WithLogFilename(string) ServerConfig
	WithLogMaxSize(int) ServerConfig
	WithLogMaxBackups(int) ServerConfig
	WithLogMaxAge(int) ServerConfig
	WithLogCompress(bool) ServerConfig
}
