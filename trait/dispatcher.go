package trait

import "time"

type Dispatcher interface {
	Start()
	Dispatch(connQueue chan Connection)
	SetHeaderDeadline(deadline time.Time)
	SetBodyDeadline(deadline time.Time)
	ChooseQueue(connID uint64) chan <- Connection
	Commit(conn Connection)
}
