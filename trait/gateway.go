package trait

type Gateway interface {
	ListenAndServe() error
	Accept() (Connection, error)
	Stop() error
}
