package log

import (
	"github.com/getto-systems/project-example-id/log"

	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (logger Logger) log() ticket_infra.Logger {
	return logger
}
