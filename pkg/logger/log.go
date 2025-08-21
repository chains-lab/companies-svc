package logger

import (
	"errors"
	"strings"

	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/sirupsen/logrus"
)

func NewLogger(cfg config.Config) *logrus.Logger {
	log := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(cfg.Server.Log.Level))
	if err != nil {
		log.Warnf("invalid log level '%s', defaulting to 'info'", cfg.Server.Log.Level)
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	switch strings.ToLower(cfg.Server.Log.Format) {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}

type Logger interface {
	WithError(err error) *logrus.Entry

	logrus.FieldLogger
}

type logger struct {
	*logrus.Entry
}

func (l *logger) WithError(err error) *logrus.Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return l.Entry.WithError(ae.Unwrap())
	}
	return l.Entry.WithError(err)
}

func NewWithBase(base *logrus.Logger) Logger {
	log := logger{
		Entry: logrus.NewEntry(base),
	}

	return &log
}
