package trait

type TaskMgr[T any] interface {
	RouterGroup[T]

	Start()
	StartWorker(taskQueue <- chan Request[T])
	ChooseQueue(connID uint64) chan <- Request[T]
	Submit(request Request[T])
}
