package rlogging

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tech-scripter/sandbox/env"
	"github.com/tech-scripter/sandbox/logging"
)

// RlogAdapter is an adapter for logrus that implements the logging.Logger interface.
type RlogAdapter struct {
	*logrus.Entry
}

// New initializes and returns a RlogAdapter configured based on environment variables.
//
// The following is a list of environment variables used:
//
//	APP_VERSION (string) sets app.version which is part of all log messages.
//	HOST (string) sets app.host which is part of all log messages.
//	LOG_JSON (bool) controls whether to output logs in JSON format with timestamps in time.RFC3339Nano format. Defaults to false.
//	LOG_CALLERS (bool) controls whether to include the calling method as a field in logs. Defaults to false. https://godoc.org/github.com/sirupsen/logrus#SetReportCaller
//	LOG_LEVEL (string) sets the log level. Defaults to "trace". https://godoc.org/github.com/sirupsen/logrus#Level
//
// If LOG_JSON is enabled, timestamps are switched
func New() *RlogAdapter {
	l := logrus.New()
	if env.GetBool("LOG_JSON") {
		l.Formatter = &logrus.JSONFormatter{
			DataKey:         "data",
			TimestampFormat: time.RFC3339Nano,
		}
	}

	l.SetReportCaller(env.GetBool("LOG_CALLERS"))

	if lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		l.Level = lvl
	} else {
		l.Level = logrus.TraceLevel
	}

	log := l.WithField("app", map[string]string{
		"host":    os.Getenv("HOST"),
		"version": os.Getenv("APP_VERSION"),
	})

	return &RlogAdapter{log}
}

func (r *RlogAdapter) With(groupName string, args ...any) logging.Logger {
	if groupName != "" {
		return &RlogAdapter{r.Entry.WithField(groupName, argsToMap(args...))}
	}

	return &RlogAdapter{r.with(args...)}
}

// WithError returns logging.Logger with the error field set.
func (r *RlogAdapter) WithError(err error) logging.Logger {
	return &RlogAdapter{r.Entry.WithError(err)}
}

// Trace logs the message at the TRACE level.
func (r *RlogAdapter) Trace(args ...interface{}) {
	r.Entry.Trace(args...)
}

// Debug logs the message at the DEBUG level.
func (r *RlogAdapter) Debug(args ...interface{}) {
	r.Entry.Debug(args...)
}

// Debugf logs a formatted debug message at the DEBUG level.
func (r *RlogAdapter) Debugf(format string, args ...any) {
	r.Entry.Debugf(format, args...)
}

// Info logs the message at the INFO level.
func (r *RlogAdapter) Info(args ...interface{}) {
	r.Entry.Info(args...)
}

// Infof logs a formatted debug message at the INFO level.
func (r *RlogAdapter) Infof(format string, args ...any) {
	r.Entry.Infof(format, args...)
}

// Error logs the message at the ERROR level.
func (r *RlogAdapter) Error(args ...interface{}) {
	r.Entry.Error(args...)
}

// Fatal logs the message at the FATAL level.
func (r *RlogAdapter) Fatal(args ...interface{}) {
	r.Entry.Fatal(args...)
}

// Log logs a message at the specified level with context.
func (r *RlogAdapter) Log(ctx context.Context, level int, msg string, args ...any) {
	entry := r.Entry.WithContext(ctx)
	entry = r.with(args...)
	entry.Log(logrus.Level(level), msg)
}

// Printf logs a message with formatting.
func (r *RlogAdapter) Printf(format string, args ...interface{}) {
	r.Entry.Printf(format, args...)
}

func (r *RlogAdapter) with(args ...any) *logrus.Entry {
	fields := argsToMap(args...)
	entry := r.Entry
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	return entry
}

func argsToMap(args ...any) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if ok {
				fields[key] = args[i+1]
			}
		}
	}
	return fields
}
