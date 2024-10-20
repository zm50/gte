package gcore

import (
	"fmt"

	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/gconf"
	"github.com/zm50/gte/glog"
	"github.com/zm50/gte/trait"
)

// Engine 服务器引擎接口实现
type Engine[T any] struct {
	trait.ServerConfig

	gateway trait.Gateway[T]
	connMgr trait.ConnMgr[T]
	taskMgr trait.TaskMgr[T]
}

// NewEngine 创建一个新的服务器引擎实例
func NewEngine[T any]() (*Engine[T], error) {
	// 新建任务管理器
	taskMgr := NewTaskMgr[T]()

	connMgr, err := NewConnMgr(gconf.Config.EpollTimeout(), gconf.Config.EpollEventSize(), taskMgr)
	if err != nil {
		glog.Error("NewConnMgr error:", err)
		return nil, err
	}

	var gateway trait.Gateway[T]
	switch gconf.Config.NetworkMode() {
	case constant.TCPNetowrkMode:
		gateway = NewTCPGateway(connMgr, taskMgr)
	case constant.WebsocketNetworkMode:
		gateway = NewWebsocketGateway(connMgr, taskMgr)
	default:
		gateway = NewTCPGateway(connMgr, taskMgr)
	}

	engine := &Engine[T]{
		ServerConfig: gconf.Config,
		gateway:      gateway,
		connMgr:      connMgr,
		taskMgr:      taskMgr,
	}

	return engine, nil
}

// Run 启动服务器引擎
func (e *Engine[T]) Run() error {
	glog.Init()

	fmt.Print(constant.Logo)
	glog.Infof("Server listening on %s:%d\n", gconf.Config.ListenIP(), gconf.Config.ListenPort())

	e.taskMgr.Start()
	go e.connMgr.Start()

	err := e.gateway.ListenAndServe()
	if err != nil {
		glog.Error("ListenAndServe error:", err)
		return err
	}

	return nil
}

// Regist 注册任务处理逻辑
func (e *Engine[T]) Regist(id uint32, flow ...TaskFunc[T]) {
	for _, fn := range flow {
		e.taskMgr.Regist(id, fn)
	}
}

// RegistFlow 注册任务处理流
func (e *Engine[T]) RegistFlow(id uint32, flow trait.TaskFlow[T]) {
	e.taskMgr.RegistFlow(id, flow)
}

// TaskFlow 获取任务处理流
func (e *Engine[T]) TaskFlow(id uint32) trait.TaskFlow[T] {
	return e.taskMgr.TaskFlow(id)
}

// Group 路由分组
func (e *Engine[T]) Group(flow ...trait.TaskFunc[T]) trait.RouterGroup[T] {
	return e.taskMgr.Group(flow...)
}

// Use 注册插件
func (e *Engine[T]) Use(flow ...trait.TaskFunc[T]) {
	e.taskMgr.Use(flow...)
}

// OnConnStart 注册连接建立的回调函数
func (e *Engine[T]) OnConnStart(fn func(conn trait.Connection[T])) {
	e.connMgr.OnConnStart(fn)
}

// OnConnStop 注册连接断开的回调函数
func (e *Engine[T]) OnConnStop(fn func(conn trait.Connection[T])) {
	e.connMgr.OnConnStop(fn)
}

// OnConnActive 注册连接变为不活跃状态的回调函数
func (e *Engine[T]) OnConnNotActive(fn func(conn trait.Connection[T])) {
	e.connMgr.OnConnNotActive(fn)
}
