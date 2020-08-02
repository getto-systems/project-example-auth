package log

import (
	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/user"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() user.Logger {
	return log
}
