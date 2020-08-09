package log

import (
	"github.com/getto-systems/project-example-id/log"

	infra "github.com/getto-systems/project-example-id/infra/credential"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (logger Logger) log() infra.Logger {
	return logger
}
