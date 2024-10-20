package trait

type Gateway[T any] interface {
	ListenAndServe() error
	Accept() (Connection[T], error)
}
