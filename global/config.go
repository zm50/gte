package global

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
