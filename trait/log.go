package trait

type Log interface {
	Info(...any)
	Error(...any)
	Warn(...any)
	Fatal(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
	Warnf(string, ...any)
	Fatalf(string, ...any)
}