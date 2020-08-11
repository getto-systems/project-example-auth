package password_reset_log

import (
	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"
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
