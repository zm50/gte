package gcore

import (
	"fmt"
	"net"

	"github.com/go75/gte/global"
	"github.com/go75/gte/trait"
)

// Engine 服务器引擎接口实现
type Engine struct {
	address net.TCPAddr
	version string

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
}

var _ trait.Engine = (*Engine)(nil)

// NewEngine 创建一个新的服务器引擎实例
func NewEngine(ip string, port int, version string) (trait.Engine, error) {
	// 新建任务管理器
	taskMgr := NewTaskMgr()
	
	connMgr, err := NewConnMgr(global.Config.EpollTimeout, global.Config.EpollEventSize)
	if err != nil {
		fmt.Println("NewConnMgr error:", err)
		return nil, err
	}

	engine := &Engine{
		address: net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		},
		version: version,
		connMgr: connMgr,
		taskMgr: taskMgr,
	}

	return engine, nil
}

// Run 启动服务器引擎
func (e *Engine) Run() error {
	fmt.Println("Server listening on ", e.address.String())
	listener, err := net.ListenTCP(e.version, &e.address)
	if err != nil {
		return err
	}

	e.connMgr.Start()

	e.acceptConn(listener)

	return nil
}

// Stop 停止服务器引擎
func (e *Engine) Stop() {
	
}

// Regist 注册任务处理逻辑
func (e *Engine) Regist(id uint16, flow ...trait.TaskFunc) {
	e.taskMgr.Regist(id, flow...)
}

// RegistFlow 注册任务处理流
func (e *Engine) RegistFlow(id uint16, flow trait.TaskFlow) {
	e.taskMgr.RegistFlow(id, flow)
}

// TaskFlow 获取任务处理流
func (e *Engine) TaskFlow(id uint16) trait.TaskFlow {
	return e.taskMgr.TaskFlow(id)
}

// Group 路由分组
func (e *Engine) Group(flow ...trait.TaskFunc) trait.RouterGroup {
	return e.taskMgr.Group(flow...)
}

// Use 注册插件
func (e *Engine) Use(flow ...trait.TaskFunc) {
	e.taskMgr.Use(flow...)
}

// acceptConn 监听阻塞客户端连接
func (e *Engine) acceptConn(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		file, err := conn.File()
		if err != nil {
			fmt.Println("File error:", err)
			continue
		}
		// 处理
		connection := NewConnection(int32(file.Fd()), conn)

		e.connMgr.Add(connection)
	}
}

