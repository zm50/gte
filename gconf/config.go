package gconf

import (
	"os"

	"github.com/pkg/errors"
	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/trait"

	"gopkg.in/yaml.v3"
)

// ServerConfig gte框架内部配置
type ServerConfig struct {
	listenIP                  string
	listenPort                int
	networkVersion            string
	readTry                   int
	writeInternal             int
	networkMode               int
	maxConns                  int
	maxPacketSize             int
	epollTimeout              int
	epollEventSize            int
	dispatcherQueues          int
	dispatcherQueueLen        int
	workersPerDispatcherQueue int
	taskQueues                int
	taskQueueLen              int
	workersPerTaskQueue       int
	websocketQueueLen         int
	connSignalQueues          int
	connSignalQueueLen        int
	workersPerConnSignalQueue int
	connShardCount            int
	healthCheckInterval       int
	logFilename               string // 日志文件存放目录
	logMaxSize                int    // 文件大小限制,单位MB
	logMaxBackups             int    // 最大保留日志文件数量
	logMaxAge                 int    // 日志文件保留天数
	logCompress               bool   // 是否压缩处理
}

var _ trait.ServerConfig = (*ServerConfig)(nil)

// Config gte框架默认配置
var Config trait.ServerConfig = &ServerConfig{
	listenIP:       "0.0.0.0",
	listenPort:     8080,
	networkVersion: "tcp4",
	readTry:       1,
	writeInternal: 100,
	networkMode:    constant.TCPNetowrkMode,

	maxConns:      1024,
	maxPacketSize: 4096,

	epollTimeout:   -1,
	epollEventSize: 128,

	dispatcherQueues:          8,
	dispatcherQueueLen:        128,
	workersPerDispatcherQueue: 2,

	taskQueues:          8,
	taskQueueLen:        128,
	workersPerTaskQueue: 4,

	websocketQueueLen: 16,

	connSignalQueues:          2,
	connSignalQueueLen:        4,
	workersPerConnSignalQueue: 2,
	connShardCount:            16,
	healthCheckInterval:       120000,

	logFilename:   "./gte.log",
	logMaxSize:    100,
	logMaxBackups: 100,
	logMaxAge:     30,
	logCompress:   false,
}

// Load 从配置文件中加载配置
func (c *ServerConfig) Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// 配置文件打开失败
		return errors.WithMessage(err, "config file open failed")
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		// 配置文件解析失败
		return errors.WithMessage(err, "config file parse failed")
	}

	return nil
}

// Export 导出配置到文件
func (c *ServerConfig) Export(filePath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		// 配置文件导出失败
		return errors.WithMessage(err, "config file export failed")
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		// 配置文件保存失败
		return errors.WithMessage(err, "config file save failed")
	}

	return nil
}

func (c *ServerConfig) ListenIP() string {
	return c.listenIP
}

func (c *ServerConfig) ListenPort() int {
	return c.listenPort
}

func (c *ServerConfig) NetworkVersion() string {
	return c.networkVersion
}

func (c *ServerConfig) ReadTry() int {
	return c.readTry
}

func (c *ServerConfig) WriteInternal() int {
	return c.writeInternal
}

func (c *ServerConfig) NetworkMode() int {
	return c.networkMode
}

func (c *ServerConfig) MaxConns() int {
	return c.maxConns
}

func (c *ServerConfig) MaxPacketSize() int {
	return c.maxPacketSize
}

func (c *ServerConfig) EpollTimeout() int {
	return c.epollTimeout
}

func (c *ServerConfig) EpollEventSize() int {
	return c.epollEventSize
}

func (c *ServerConfig) DispatcherQueues() int {
	return c.dispatcherQueues
}

func (c *ServerConfig) DispatcherQueueLen() int {
	return c.dispatcherQueueLen
}

func (c *ServerConfig) WorkersPerDispatcherQueue() int {
	return c.workersPerDispatcherQueue
}

func (c *ServerConfig) TaskQueues() int {
	return c.taskQueues
}

func (c *ServerConfig) TaskQueueLen() int {
	return c.taskQueueLen
}

func (c *ServerConfig) WorkersPerTaskQueue() int {
	return c.workersPerTaskQueue
}

func (c *ServerConfig) WebsocketQueueLen() int {
	return c.websocketQueueLen
}

func (c *ServerConfig) ConnSignalQueues() int {
	return c.connSignalQueues
}

func (c *ServerConfig) ConnSignalQueueLen() int {
	return c.connSignalQueueLen
}

func (c *ServerConfig) WorkersPerConnSignalQueue() int {
	return c.workersPerConnSignalQueue
}

func (c *ServerConfig) ConnShardCount() int {
	return c.connShardCount
}

func (c *ServerConfig) HealthCheckInterval() int {
	return c.healthCheckInterval
}

func (c *ServerConfig) LogFilename() string {
	return c.logFilename
}

func (c *ServerConfig) LogMaxSize() int {
	return c.logMaxSize
}

func (c *ServerConfig) LogMaxBackups() int {
	return c.logMaxBackups
}

func (c *ServerConfig) LogMaxAge() int {
	return c.logMaxAge
}

func (c *ServerConfig) LogCompress() bool {
	return c.logCompress
}

func (c *ServerConfig) WithListenIP(listenIP string) trait.ServerConfig {
	c.listenIP = listenIP
	return c
}

func (c *ServerConfig) WithListenPort(listenPort int) trait.ServerConfig {
	c.listenPort = listenPort
	return c
}

func (c *ServerConfig) WithNetworkVersion(networkVersion string) trait.ServerConfig {
	c.networkVersion = networkVersion
	return c
}

func (c *ServerConfig) WithReadTry(readTry int) trait.ServerConfig {
	c.readTry = readTry
	return c
}

func (c *ServerConfig) WithWriteInternal(writeInternal int) trait.ServerConfig {
	c.writeInternal = writeInternal
	return c
}

func (c *ServerConfig) WithNetworkMode(networkMode int) trait.ServerConfig {
	c.networkMode = networkMode
	return c
}

func (c *ServerConfig) WithMaxConns(maxConns int) trait.ServerConfig {
	c.maxConns = maxConns
	return c
}

func (c *ServerConfig) WithMaxPacketSize(maxPacketSize int) trait.ServerConfig {
	c.maxPacketSize = maxPacketSize
	return c
}

func (c *ServerConfig) WithEpollTimeout(epollTimeout int) trait.ServerConfig {
	c.epollTimeout = epollTimeout
	return c
}

func (c *ServerConfig) WithEpollEventSize(epollEventSize int) trait.ServerConfig {
	c.epollEventSize = epollEventSize
	return c
}

func (c *ServerConfig) WithDispatcherQueues(dispatcherQueues int) trait.ServerConfig {
	c.dispatcherQueues = dispatcherQueues
	return c
}

func (c *ServerConfig) WithDispatcherQueueLen(dispatcherQueueLen int) trait.ServerConfig {
	c.dispatcherQueueLen = dispatcherQueueLen
	return c
}

func (c *ServerConfig) WithWorkersPerDispatcherQueue(workersPerDispatcherQueue int) trait.ServerConfig {
	c.workersPerDispatcherQueue = workersPerDispatcherQueue
	return c
}

func (c *ServerConfig) WithTaskQueues(taskQueues int) trait.ServerConfig {
	c.taskQueues = taskQueues
	return c
}

func (c *ServerConfig) WithTaskQueueLen(taskQueueLen int) trait.ServerConfig {
	c.taskQueueLen = taskQueueLen
	return c
}

func (c *ServerConfig) WithWorkersPerTaskQueue(workersPerTaskQueue int) trait.ServerConfig {
	c.workersPerTaskQueue = workersPerTaskQueue
	return c
}

func (c *ServerConfig) WithWebsocketQueueLen(websocketQueueLen int) trait.ServerConfig {
	c.websocketQueueLen = websocketQueueLen
	return c
}

func (c *ServerConfig) WithConnSignalQueues(connSignalQueues int) trait.ServerConfig {
	c.connSignalQueues = connSignalQueues
	return c
}

func (c *ServerConfig) WithConnSignalQueueLen(connSignalQueueLen int) trait.ServerConfig {
	c.connSignalQueueLen = connSignalQueueLen
	return c
}

func (c *ServerConfig) WithWorkersPerConnSignalQueue(workersPerConnSignalQueue int) trait.ServerConfig {
	c.workersPerConnSignalQueue = workersPerConnSignalQueue
	return c
}

func (c *ServerConfig) WithConnShardCount(connShardCount int) trait.ServerConfig {
	c.connShardCount = connShardCount
	return c
}

func (c *ServerConfig) WithHealthCheckInterval(healthCheckInterval int) trait.ServerConfig {
	c.healthCheckInterval = healthCheckInterval
	return c
}

func (c *ServerConfig) WithLogFilename(logFilename string) trait.ServerConfig {
	c.logFilename = logFilename
	return c
}

func (c *ServerConfig) WithLogMaxSize(logMaxSize int) trait.ServerConfig {
	c.logMaxSize = logMaxSize
	return c
}

func (c *ServerConfig) WithLogMaxBackups(logMaxBackups int) trait.ServerConfig {
	c.logMaxBackups = logMaxBackups
	return c
}

func (c *ServerConfig) WithLogMaxAge(logMaxAge int) trait.ServerConfig {
	c.logMaxAge = logMaxAge
	return c
}

func (c *ServerConfig) WithLogCompress(logCompress bool) trait.ServerConfig {
	c.logCompress = logCompress
	return c
}
