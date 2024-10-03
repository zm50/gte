package trait

type RouterGroup interface {
	Router

	Group(flow ...TaskFunc) RouterGroup
	Use(flow ...TaskFunc)
}
