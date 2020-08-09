package user_log

import (
	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/user/infra"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() infra.Logger {
	return log
}
