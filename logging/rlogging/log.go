package rlogging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tech-scripter/sandbox/env"
)

type RLogger struct {
	*logrus.Entry
}

func New() *RLogger {
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

	logger := l.WithFields(logrus.Fields{
		"app": map[string]string{
			"host":    os.Getenv("HOST"),
			"version": os.Getenv("APP_VERSION"),
		},
	})

	return &RLogger{logger}
}

func (r *RLogger) WithError(err error) logging.Logger {
	return &RLogger{r.Logger.WithError(err)}
}

func (r *RLogger) Error(args ...interface{}) {
	r.Logger.Error(args)
}

func (r *RLogger) Printf(format string, args ...interface{}) {
	r.Logger.Printf(format, args)
}
