package log

import (
	"github.com/getto-systems/project-example-id/log"

	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (logger Logger) log() credential_infra.Logger {
	return logger
}
