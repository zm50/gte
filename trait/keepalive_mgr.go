package trait

import "github.com/zm50/gte/core"

type KeepAliveMgr[T any] interface {
	Start()
	StartWorker(connShard *core.KVShard[int32, Connection[T]])
}
