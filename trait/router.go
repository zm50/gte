package trait

type Router interface {
	Regist(id uint32, flow ...TaskFunc)
	RegistFlow(id uint32, flow TaskFlow)
	TaskFlow(id uint32) TaskFlow
}
