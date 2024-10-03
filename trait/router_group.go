package trait

type RouterGroup interface {
	Group(flow ...TaskFunc) RouterGroup
	Use(flow ...TaskFunc)
	Regist(id uint16, flow ...TaskFunc)
	RegistFlow(id uint16, flow TaskFlow)
}
