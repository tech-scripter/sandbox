package slogging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/tech-scripter/sandbox/env"
	"github.com/tech-scripter/sandbox/logging"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(9)
)

var (
	LevelNames = map[slog.Leveler]string{
		LevelTrace: "TRACE",
		LevelFatal: "FATAL",
	}

	exit = func(code int) {
		os.Exit(code)
	}
)

// SlogAdapter is an adapter for slog.Logger that implements the logging.Logger interface.
type SlogAdapter struct {
	*slog.Logger
}

// New initializes and returns a SlogAdapter configured based on environment variables.
//
// The following is a list of environment variables used:
//
//	APP_VERSION (string) sets app.version which is part of all log messages.
//	HOST (string) sets app.host which is part of all log messages.
//	LOG_JSON (bool) controls whether to output logs in JSON format with timestamps in time.RFC3339Nano format. Defaults to false.
//	LOG_CALLERS (bool) controls whether to include the calling method as a field in logs. Defaults to false. https://pkg.go.dev/log/slog#HandlerOptions
//	LOG_LEVEL (string) sets the log level. Defaults to "info". https://pkg.go.dev/log/slog#Level
//
// If LOG_JSON is enabled, timestamps are switched
func New() *SlogAdapter {
	handlerOptions := &slog.HandlerOptions{
		Level: LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}

	if env.GetBool("LOG_CALLERS") {
		handlerOptions.AddSource = true
	}

	if level, err := parseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		handlerOptions.Level = level
	}

	var handler slog.Handler
	if env.GetBool("LOG_JSON") {
		handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
	} else {
		handler = slog.NewTextHandler(os.Stderr, handlerOptions)
	}

	logger := slog.New(handler).WithGroup("data").With(
		slog.Group("app",
			slog.String("host", os.Getenv("HOST")),
			slog.String("version", os.Getenv("APP_VERSION")),
		),
	)

	return &SlogAdapter{logger}
}

// With adds key-value pairs to the logger and returns logging.Logger.
func (s *SlogAdapter) With(groupName string, args ...any) logging.Logger {
	if groupName != "" {
		return &SlogAdapter{s.Logger.WithGroup(groupName).With(args...)}
	}

	return &SlogAdapter{s.Logger.With(args...)}
}

// WithError automatically adds the error key and returns logging.Logger
func (s *SlogAdapter) WithError(err error) logging.Logger {
	return &SlogAdapter{s.Logger.With("error", err)}
}

// Trace logs the message at the TRACE level.
func (s *SlogAdapter) Trace(args ...interface{}) {
	s.Logger.Log(context.Background(), LevelTrace, fmt.Sprint(args...))
}

// Debug logs the message at the DEBUG level.
func (s *SlogAdapter) Debug(args ...interface{}) {
	s.Logger.Debug(fmt.Sprint(args...))
}

// Debugf logs the message at the DEBUG level using a formatted string.
func (s *SlogAdapter) Debugf(format string, args ...any) {
	s.Logger.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
}

// Info logs the message at the INFO level.
func (s *SlogAdapter) Info(args ...interface{}) {
	s.Logger.Info(fmt.Sprint(args...))
}

// Infof logs the message at the INFO level using a formatted string.
func (s *SlogAdapter) Infof(format string, args ...any) {
	s.Logger.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

// Error logs the message at the ERROR level.
func (s *SlogAdapter) Error(args ...interface{}) {
	s.Logger.Error(fmt.Sprint(args...))
}

// Fatal logs the message and then exits the program.
func (s *SlogAdapter) Fatal(args ...interface{}) {
	s.Logger.Log(context.Background(), LevelFatal, fmt.Sprint(args...))
	exit(1)
}

// Log logs a message at the specified level with context.
func (s *SlogAdapter) Log(ctx context.Context, level int, msg string, args ...any) {
	s.Logger.With(args...).Log(ctx, slog.Level(level), msg)
}

// Printf logs the message at the INFO level using a formatted string.
func (s *SlogAdapter) Printf(format string, args ...interface{}) {
	s.Infof(format, args...)
}

// parseLevel parses a string into a slog.Level.
// It returns an error and sets the level to INFO if the string is not a valid level.
func parseLevel(level string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(level))
	return l, err
}
