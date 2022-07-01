package log

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

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
