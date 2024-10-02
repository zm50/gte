package trait

type TaskMgr interface {
	Start()
	StartWorker(taskQueue chan Request)
	Submit(request Request)
	Regist(id uint16, flow TaskFlow)
}
