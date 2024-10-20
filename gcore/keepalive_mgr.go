package gcore

import (
	"time"

	"github.com/zm50/gte/constant"
	"github.com/zm50/gte/core"
	"github.com/zm50/gte/gconf"
	"github.com/zm50/gte/glog"
	"github.com/zm50/gte/trait"
)

// KeepAliveMgr 连接存活管理器
type KeepAliveMgr[T any] struct {
	connMgr             trait.ConnMgr[T]
	healthCheckInterval time.Duration
	connShards          []*core.KVShard[int32, trait.Connection[T]]
}

// NewKeepAliveMgr 创建连接存活管理器
func NewKeepAliveMgr[T any](connMgr trait.ConnMgr[T], connShards []*core.KVShard[int32, trait.Connection[T]]) trait.KeepAliveMgr[T] {
	return &KeepAliveMgr[T]{
		connMgr:             connMgr,
		healthCheckInterval: time.Millisecond * time.Duration(gconf.Config.HealthCheckInterval()),
		connShards:          connShards,
	}
}

// NewKeepAliveMgr 启动连接存活管理器
func (m *KeepAliveMgr[T]) Start() {
	glog.Info("keepalive manager start...")

	for _, connShard := range m.connShards {
		go m.StartWorker(connShard)
	}
}

// StartWorker 启动健康检查工作
func (k *KeepAliveMgr[T]) StartWorker(connShard *core.KVShard[int32, trait.Connection[T]]) {
	ticker := time.NewTicker(k.healthCheckInterval)
	for {
		<-ticker.C
		connShard.RRange(func(id int32, conn trait.Connection[T]) {
			state := conn.State()
			if state == constant.ConnActiveState {
				// 设置为检查状态
				conn.SetState(constant.ConnInspectState)
			} else if state == constant.ConnInspectState {
				// 设置为非活跃状态
				conn.SetState(constant.ConnNotActiveState)
				k.connMgr.PushConnSignal(NewConnSignal(conn, constant.ConnNotActiveSignal))
			}
		})
	}
}
