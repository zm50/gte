package trait

type Router interface {
	Regist(id uint16, flow ...TaskFunc)
	RegistFlow(id uint16, flow TaskFlow)
	TaskFlow(id uint16) TaskFlow
}
