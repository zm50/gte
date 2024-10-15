package gcore

import (
	"time"

	"github.com/go75/gte/constant"
	"github.com/go75/gte/core"
	"github.com/go75/gte/gconf"
	"github.com/go75/gte/trait"
)

// KeepAliveMgr 连接存活管理器
type KeepAliveMgr struct {
	connMgr trait.ConnMgr
	healthCheckInterval time.Duration
	connShards []*core.KVShard[int32, trait.Connection]
}

// NewKeepAliveMgr 创建连接存活管理器
func NewKeepAliveMgr(connMgr trait.ConnMgr, connShards []*core.KVShard[int32, trait.Connection]) trait.KeepAliveMgr {
	return &KeepAliveMgr{
		connMgr: connMgr,
		healthCheckInterval: time.Millisecond * time.Duration(gconf.Config.HealthCheckInterval()),
		connShards: connShards,
	}
}

// NewKeepAliveMgr 启动连接存活管理器
func (m *KeepAliveMgr) Start() {
	for _, connShard := range m.connShards {
		go m.StartWorker(connShard)
	}
}

// StartWorker 启动健康检查工作
func (k *KeepAliveMgr) StartWorker(connShard *core.KVShard[int32, trait.Connection]) {
	ticker := time.NewTicker(k.healthCheckInterval)
	for {
		<- ticker.C
		connShard.RRange(func(id int32, conn trait.Connection) {
			if conn.IsActive() {
				// 设置为检查状态
				conn.SetState(constant.ConnInspectState)
			} else if conn.IsInspect() {
				// 设置为非活跃状态
				conn.SetState(constant.ConnNotActiveState)
				k.connMgr.PushConnSignal(NewConnSignal(conn, constant.ConnNotActiveSignal))
			}
		})
	}
}
