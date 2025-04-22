package logging

type Logger interface {
	Printf(format string, args ...interface{})
	WithError(err error) Logger
	Error(args ...interface{})
}
