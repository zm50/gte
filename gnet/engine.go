package gnet

import (
	"fmt"
	"net"

	"github.com/go75/gte/global"
	"github.com/go75/gte/trait"
)

// Engine 服务器引擎接口实现
type Engine struct {
	trait.RouterGroup
	address net.TCPAddr
	version string

	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
}

var _ trait.Engine = (*Engine)(nil)

// NewEngine 创建一个新的服务器引擎实例
func NewEngine(ip string, port int, version string) (trait.Engine, error) {
	// 新建任务处理路由器
	router := NewRouter()

	// 新建任务管理器
	taskMgr := NewTaskMgr(router)
	
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

	engine.RouterGroup = NewRouterGroup(engine)

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

