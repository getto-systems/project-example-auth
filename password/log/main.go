package log

import (
	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() password.Logger {
	return log
}
