package handler

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tech-scripter/sandbox/env"
)

var (
	log *logrus.Entry

	logrusLevelsBySlogLevels = map[slog.Level]logrus.Level{
		slog.LevelDebug: logrus.DebugLevel,
		slog.LevelInfo:  logrus.InfoLevel,
		slog.LevelWarn:  logrus.WarnLevel,
		slog.LevelError: logrus.ErrorLevel,
	}
)

type LogrusHandler struct {
	Logger *logrus.Entry
}

func New() *LogrusHandler {
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

	log = l.WithFields(logrus.Fields{
		"app": map[string]string{
			"host":    os.Getenv("HOST"),
			"version": os.Getenv("APP_VERSION"),
		},
	})

	return &LogrusHandler{Logger: log}
}

func (l *LogrusHandler) Enabled(_ context.Context, level slog.Level) bool {
	logrusLevel := logrus.TraceLevel // default level
	if lvl, exists := logrusLevelsBySlogLevels[level]; exists {
		logrusLevel = lvl
	}
	return l.Logger.Logger.IsLevelEnabled(logrusLevel)
}

func (l *LogrusHandler) Handle(ctx context.Context, rec slog.Record) error {
	fields := make(map[string]interface{}, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	entry := l.Logger.WithContext(ctx).WithFields(fields)

	switch rec.Level {
	case slog.LevelDebug:
		entry.Debug(rec.Message)
	case slog.LevelInfo:
		entry.Info(rec.Message)
	case slog.LevelWarn:
		entry.Warn(rec.Message)
	case slog.LevelError:
		entry.Error(rec.Message)
	default:
		entry.Trace(rec.Message)
	}

	return nil
}

func (l *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make(map[string]interface{}, len(attrs))
	for _, attr := range attrs {
		fields[attr.Key] = attr.Value.Any()
	}
	return &LogrusHandler{Logger: l.Logger.WithFields(fields)}
}

func (l *LogrusHandler) WithGroup(name string) slog.Handler {
	return &LogrusHandler{Logger: l.Logger.WithField("group", name)}
}
