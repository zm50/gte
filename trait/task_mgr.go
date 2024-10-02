package trait

import "time"

type TaskMgr interface {
	Start()
	StartWorker(taskQueue chan Request)
	Submit(request Request)
	Regist(id uint16, flow TaskFlow)
	BatchDispatch(conn Connection, headerDeadline, bodyDeadline time.Time) error
}
