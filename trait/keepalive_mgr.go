package trait

type KeepAliveMgr[T any] interface {
	Start()
	StartWorker(connShard KVShard[int32, Connection[T]])
}
