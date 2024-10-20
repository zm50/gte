package trait

import "time"

type Dispatcher[T any] interface {
	Start()
	Dispatch(connQueue chan Connection[T])
	SetHeaderDeadline(deadline time.Time)
	SetBodyDeadline(deadline time.Time)
	ChooseQueue(connID uint64) chan <- Connection[T]
	Commit(conn Connection[T])
}
