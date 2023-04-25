package httpclient

import "github.com/sirupsen/logrus"

type Logger interface {
	Error(args ...any)
}

func DefaultLog() Logger {
	return logrus.StandardLogger()
}
