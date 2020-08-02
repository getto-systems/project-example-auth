package log

import (
	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/api_token"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (logger Logger) log() api_token.Logger {
	return logger
}
