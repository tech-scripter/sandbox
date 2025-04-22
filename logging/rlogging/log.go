package rlogging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tech-scripter/sandbox/env"
	"github.com/tech-scripter/sandbox/logging"
)

type RlogAdapter struct {
	*logrus.Entry
}

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

	logger := l.WithFields(logrus.Fields{
		"app": map[string]string{
			"host":    os.Getenv("HOST"),
			"version": os.Getenv("APP_VERSION"),
		},
	})

	return &RlogAdapter{logger}
}

func (r *RlogAdapter) WithError(err error) logging.Logger {
	return &RlogAdapter{r.Logger.WithError(err)}
}

func (r *RlogAdapter) Error(args ...interface{}) {
	r.Logger.Error(args)
}

func (r *RlogAdapter) Printf(format string, args ...interface{}) {
	r.Logger.Printf(format, args)
}
