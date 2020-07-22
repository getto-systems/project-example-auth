package log

import (
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

type Logger struct {
	logger event_log.Logger
}

func NewLogger(logger event_log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() password.Logger {
	return log
}
