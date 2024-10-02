package trait

import "time"

type Dispatcher interface {
	Start()
	Dispatch(connQueue chan Connection)
	BatchDispatch(conn Connection) error
	SetHeaderDeadline(deadline time.Time)
	SetBodyDeadline(deadline time.Time)
	ChooseQueue(conn Connection) chan <- Connection
	Commit(conn Connection)
}
