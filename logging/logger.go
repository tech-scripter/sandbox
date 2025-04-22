package logging

// Logger is an interface that wraps the basic Printf method.
// This interface provides a unified way to use loggers from various third-party libraries as function parameters.
// For instance, the SetTraceLog() in github.com/olivere/elastic/v7 accepts a logger of type Logger.
type Logger interface {
	Printf(format string, args ...interface{})
	WithError(err error) Logger
	Error(args ...interface{})
}
