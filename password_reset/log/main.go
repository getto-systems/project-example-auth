package log

import (
	"github.com/getto-systems/project-example-id/log"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() password_reset_infra.Logger {
	return log
}
