package trait

import (
	"sync"
)

type ConnMgr[T any] interface {
	Get(fd int32) (Connection[T], bool)
	Add(conn Connection[T]) error
	Del(fd int32) error
	Wait() (int, error)
	BatchCommit(n int)
	Start()
	Stop()
	StartConnSignalHookWorkers()
	StartConnSignalHookWorker(<- chan ConnSignal[T])
	OnConnStart(func(conn Connection[T]))
	OnConnStop(func(conn Connection[T]))
	OnConnNotActive(fn func(conn Connection[T]))
	ChooseConnSignalQueue(connID uint64) chan <- ConnSignal[T]
	PushConnSignal(signal ConnSignal[T])
	WaitGroup() *sync.WaitGroup
	OnlineConns() int32
}
