package logging

import "context"

// Logger is an interface that defines the methods for logging.
type Logger interface {
	With(groupName string, args ...any) Logger
	WithError(err error) Logger

	Trace(args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...any)

	Info(args ...interface{})
	Infof(format string, args ...any)

	Error(args ...interface{})

	Fatal(args ...interface{})

	Log(ctx context.Context, level int, msg string, args ...any)

	Printf(format string, args ...interface{})
}
