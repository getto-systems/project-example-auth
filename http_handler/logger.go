package http_handler

import (
	"github.com/getto-systems/project-example-id/data"
)

type Logger struct {
	logger  RequestLogger
	request data.Request
}

func NewLogger(logger RequestLogger, request data.Request) Logger {
	return Logger{
		logger:  logger,
		request: request,
	}
}

type RequestLogger interface {
	DebugMessage(data.Request, string)
	DebugError(data.Request, string, error)
}

func (logger Logger) DebugMessage(message string) {
	logger.logger.DebugMessage(logger.request, message)
}

func (logger Logger) DebugError(format string, err error) {
	logger.logger.DebugError(logger.request, format, err)
}
