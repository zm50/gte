package gcore

import (
	"fmt"

	"github.com/go75/gte/constant"
	"github.com/go75/gte/gconf"
	"github.com/go75/gte/trait"
)

// Engine 服务器引擎接口实现
type Engine struct {
	trait.ServerConfig

	gateway trait.Gateway
	connMgr trait.ConnMgr
	taskMgr trait.TaskMgr
}

var _ trait.Engine = (*Engine)(nil)

// NewEngine 创建一个新的服务器引擎实例
func NewEngine(ip string, port int, version string) (trait.Engine, error) {
	// 新建任务管理器
	taskMgr := NewTaskMgr()
	
	connMgr, err := NewConnMgr(gconf.Config.EpollTimeout(), gconf.Config.EpollEventSize())
	if err != nil {
		fmt.Println("NewConnMgr error:", err)
		return nil, err
	}

	var gateway trait.Gateway
	switch gconf.Config.NetworkMode() {
	case constant.TCPNetowrkMode:
		gateway = NewTCPGateway(connMgr, taskMgr)
	case constant.WebsocketNetworkMode:
		gateway = NewWebsocketGateway(connMgr, taskMgr)
	default:
		gateway = NewTCPGateway(connMgr, taskMgr)
	}

	engine := &Engine{
		ServerConfig: gconf.Config,
		gateway: gateway,
		connMgr: connMgr,
		taskMgr: taskMgr,
	}

	return engine, nil
}

// Run 启动服务器引擎
func (e *Engine) Run() error {
	fmt.Printf("Server listening on %s:%d\n", gconf.Config.ListenIP(), gconf.Config.ListenPort())

	go e.connMgr.Start()

	err := e.gateway.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe error:", err)
		return err
	}

	return nil
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

// OnConnStart 注册连接建立的回调函数
func (e *Engine) OnConnStart(fn func(conn trait.Connection)) {
	e.connMgr.OnConnStart(fn)
}

// OnConnStop 注册连接断开的回调函数
func (e *Engine) OnConnStop(fn func(conn trait.Connection)) {
	e.connMgr.OnConnStop(fn)
}
