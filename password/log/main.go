package log

import (
	"github.com/getto-systems/project-example-id/log"

	password_infra "github.com/getto-systems/project-example-id/infra/password"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() password_infra.Logger {
	return log
}
