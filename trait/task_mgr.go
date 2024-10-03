package trait

type TaskMgr interface {
	RouterGroup

	Start()
	StartWorker(taskQueue chan Request)
	Submit(request Request)
}
