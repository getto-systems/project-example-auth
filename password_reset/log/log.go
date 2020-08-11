package password_reset_log

import (
	"github.com/getto-systems/project-example-id/password_reset/infra"
)

type (
	LeveledLogger interface {
		Audit(interface{})
		Error(interface{})
		Info(interface{})
		Debug(interface{})
	}

	Logger struct {
		logger LeveledLogger
	}
)

func NewLogger(logger LeveledLogger) Logger {
	return Logger{
		logger: logger,
	}
}

func (log Logger) handler() infra.Logger {
	return log
}
