package trait

type Dispatcher interface {
	Start()
	Dispatch(connQueue chan Connection)
	BatchDispatch(conn Connection) error
}
