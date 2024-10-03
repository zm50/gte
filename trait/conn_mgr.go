package trait

type ConnMgr interface {
	Get(fd int32) (Connection, bool)
	Add(conn Connection) error
	Del(fd int) error
	Wait() (int, error)
	BatchCommit(n int)
	Start()
	Stop()
}
