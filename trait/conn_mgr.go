package trait

import "time"

type ConnMgr interface {
	Get(fd int32) (Connection, bool)
	Add(conn Connection) error
	Del(fd int32) error
	Wait() (int, error)
	BatchCommit(n int)
	Start()
	Stop()
	ReadDeadline() time.Time
	MaxReadDeadline() time.Time
	StartConnSignalHookWorkers()
	StartConnSignalHookWorker(<- chan ConnSignal)
	OnConnStart(func(conn Connection))
	OnConnStop(func(conn Connection))
	OnConnNotActive(fn func(conn Connection))
	ChooseConnSignalQueue(connID uint64) chan <- ConnSignal
	PushConnSignal(signal ConnSignal)
}
