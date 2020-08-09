package credential_log

import (
	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/credential/infra"
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
