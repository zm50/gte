package trait

type ObjPool[T any] interface {
	Get() T
	Put(T)
}
