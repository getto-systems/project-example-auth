package ticket_log

import (
	"github.com/getto-systems/project-example-id/ticket/infra"
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

func (logger Logger) log() infra.Logger {
	return logger
}
