package trait

type TaskMgr interface {
	RouterGroup

	Start()
	StartWorker(taskQueue <- chan Request)
	ChooseQueue(connID uint64) chan <- Request
	Submit(request Request)
}
