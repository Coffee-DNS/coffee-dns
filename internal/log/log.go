package log

import (
	"github.com/sirupsen/logrus"
)

// Logger is a logger used by Coffee DNS services
type Logger struct {
	*logrus.Logger
}

// NewJSONLogger returns a new structured logger
func NewJSONLogger(level string) *Logger {
	var l Logger
	l.Logger = logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetReportCaller(true)

	// Silently use the default log level if parsing fails
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return &l
	}
	l.SetLevel(logLevel)
	return &l
}
