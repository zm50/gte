package trait

type RouterGroup[T any] interface {
	Router[T]

	Group(flow ...TaskFunc[T]) RouterGroup[T]
	Use(flow ...TaskFunc[T])
}
