package log

import (
	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/ticket"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (logger Logger) log() ticket.Logger {
	return logger
}
