package gcore

import (
	"fmt"

	"github.com/zm50/gte/common/constant"
	"github.com/zm50/gte/global"
	"github.com/zm50/gte/gpack"
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

	connMgr, err := NewConnMgr(global.Config().EpollTimeout(), global.Config().EpollEventSize(), taskMgr)
	if err != nil {
		global.Logger().Error("NewConnMgr error:", err)
		return nil, err
	}

	var gateway trait.Gateway[T]
	switch global.Config().NetworkMode() {
	case constant.TCPNetowrkMode:
		gateway = NewTCPGateway(connMgr, taskMgr)
	case constant.WebsocketNetworkMode:
		gateway = NewWebsocketGateway(connMgr, taskMgr)
	default:
		gateway = NewTCPGateway(connMgr, taskMgr)
	}

	engine := &Engine[T]{
		ServerConfig: global.Config(),
		gateway:      gateway,
		connMgr:      connMgr,
		taskMgr:      taskMgr,
	}

	return engine, nil
}

// Run 启动服务器引擎
func (e *Engine[T]) Run() error {
	e.setup()

	fmt.Print(constant.Logo)
	global.Logger().Infof("Server listening on %s:%d\n", global.Config().ListenIP(), global.Config().ListenPort())

	e.taskMgr.Start()
	go e.connMgr.Start()

	err := e.gateway.ListenAndServe()
	if err != nil {
		global.Logger().Error("ListenAndServe error:", err)
		return err
	}

	return nil
}

func (e *Engine[T]) setup() {
	global.SetLogger(global.Config().LogFilename(), global.Config().LogMaxAge(), global.Config().LogMaxBackups(),
		global.Config().LogMaxSize(), global.Config().LogCompress())

	global.SetMsgPool(global.Config().MessagePoolSize(), func() trait.Message {
		return gpack.NewMessage(0, []byte{})
	})
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
func (e *Engine[T]) Group(flow ...TaskFunc[T]) trait.RouterGroup[T] {
	fw := make([]trait.TaskFunc[T], 0, len(flow))
	for _, fn := range flow {
		fw = append(fw, fn)
	}

	return e.taskMgr.Group(fw...)
}

// Use 注册插件
func (e *Engine[T]) Use(flow ...TaskFunc[T]) {
	for _, fn := range flow {
		e.taskMgr.Use(fn)
	}
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
