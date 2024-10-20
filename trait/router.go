package trait

type Router[T any] interface {
	Regist(id uint32, flow ...TaskFunc[T])
	RegistFlow(id uint32, flow TaskFlow[T])
	TaskFlow(id uint32) TaskFlow[T]
}
