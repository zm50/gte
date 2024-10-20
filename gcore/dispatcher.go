package gcore

import (
	"time"

	"github.com/zm50/gte/gconf"
	"github.com/zm50/gte/glog"
	"github.com/zm50/gte/trait"
)

// Dispatcher 请求分发模块，负责读取客户端连接的数据，并对数据进行拆包转换成消息格式，然后分发给下游的任务处理模块对消息进行业务处理
type Dispatcher[T any] struct {
	headerDeadline time.Time
	bodyDeadline   time.Time

	connQueue []chan trait.Connection[T]
	connMgr   trait.ConnMgr[T]
	taskMgr   trait.TaskMgr[T]
}

var _ trait.Dispatcher[any] = (*Dispatcher[any])(nil)

// NewDispatcher 创建一个请求分发器
func NewDispatcher[T any](connMgr trait.ConnMgr[T], taskMgr trait.TaskMgr[T]) trait.Dispatcher[T] {
	connQueue := make([]chan trait.Connection[T], gconf.Config.DispatcherQueues())
	for i := 0; i < len(connQueue); i++ {
		connQueue[i] = make(chan trait.Connection[T], gconf.Config.DispatcherQueueLen())
	}

	return &Dispatcher[T]{
		connQueue: connQueue,
		connMgr:   connMgr,
		taskMgr:   taskMgr,
	}
}

// Start 启动请求分发模块
func (d *Dispatcher[T]) Start() {
	glog.Info("dispatcher start...")

	for i := 0; i < len(d.connQueue); i++ {
		for j := 0; j < gconf.Config.DispatcherQueueLen(); j++ {
			go d.Dispatch(d.connQueue[i])
		}
	}
}

// StartDispatcher 分发连接数据
func (d *Dispatcher[T]) Dispatch(connQueue chan trait.Connection[T]) {
	// 从conn中读取数据，并将数据提交给taskMgr处理
	for conn := range connQueue {
		err := conn.BatchCommit()
		if err != nil {
			glog.Error("dispatcher batch commit error: ", err)
			if d.connMgr.Del(int32(conn.ID())) != nil {
				glog.Error("del conn error: ", err)
			}
		}
	}
}

// SetHeaderDeadline 设置header读取超时时间
func (d *Dispatcher[T]) SetHeaderDeadline(deadline time.Time) {
	d.headerDeadline = deadline
}

// SetBodyDeadline 设置body读取超时时间
func (d *Dispatcher[T]) SetBodyDeadline(deadline time.Time) {
	d.bodyDeadline = deadline
}

// ChooseQueue 选择处理连接的队列
func (d *Dispatcher[T]) ChooseQueue(connID uint64) chan<- trait.Connection[T] {
	// 负载均衡，选择队列
	return d.connQueue[connID%uint64(len(d.connQueue))]
}

// Commit 提交连接到队列
func (d *Dispatcher[T]) Commit(conn trait.Connection[T]) {
	d.ChooseQueue(conn.ID()) <- conn
}
