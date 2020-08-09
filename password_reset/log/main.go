package log

import (
	"github.com/getto-systems/project-example-id/log"

	infra "github.com/getto-systems/project-example-id/infra/password_reset"
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
