package trait

import "github.com/go75/gte/core"

type KeepAliveMgr interface {
	Start()
	StartWorker(connShard *core.KVShard[int32, Connection])
}
