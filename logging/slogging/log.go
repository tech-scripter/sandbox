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

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

// SlogAdapter wraps slog.Logger to provide additional functionality.
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

// WithError automatically adds the error key and returns log.Logger
func (c *SlogAdapter) WithError(err error) logging.Logger {
	return &SlogAdapter{c.With("error", err)}
}

// Fatal logs the message and then exits the program.
func (c *SlogAdapter) Fatal(msg string, args ...any) {
	c.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// Trace logs the message at the TRACE level.
func (c *SlogAdapter) Trace(msg string, args ...any) {
	c.Log(context.Background(), LevelTrace, msg, args...)
}

// Debugf logs the message at the DEBUG level using a formatted string.
func (c *SlogAdapter) Debugf(format string, args ...any) {
	if c.Enabled(context.Background(), slog.LevelDebug) {
		c.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
	}
}

// Infof logs the message at the INFO level using a formatted string.
func (c *SlogAdapter) Infof(format string, args ...any) {
	if c.Enabled(context.Background(), slog.LevelInfo) {
		c.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
	}
}

// Error logs the message at the ERROR level.
func (c *SlogAdapter) Error(args ...interface{}) {
	if c.Enabled(context.Background(), slog.LevelError) {
		c.Log(context.Background(), slog.LevelError, fmt.Sprint(args...))
	}
}

// Printf logs the message at the INFO level using a formatted string.
func (c *SlogAdapter) Printf(format string, args ...interface{}) {
	c.Infof(format, args...)
}

// parseLevel parses a string into a slog.Level.
// It returns an error and sets the level to INFO if the string is not a valid level.
func parseLevel(level string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(level))
	return l, err
}
