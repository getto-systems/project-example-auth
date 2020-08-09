package password_log

import (
	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/password/infra"
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
