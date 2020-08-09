package log

import (
	"github.com/getto-systems/project-example-id/log"

	user_infra "github.com/getto-systems/project-example-id/infra/user"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() user_infra.Logger {
	return log
}
