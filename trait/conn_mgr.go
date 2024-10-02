package trait

type ConnMgr interface {
	Get(fd int32) Connection
	Add(conn Connection) error
	Del(fd int) error
	Wait() (int, error)
	BatchProcess(n int)
	Start()
	Close()
}
