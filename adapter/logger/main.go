package logger

import (
	"github.com/google/uuid"
	"log"
	"net/http"

	"github.com/getto-systems/project-example-id/misc/jsonlog"
)

type RequestLogEntry struct {
	RequestID string
	RemoteIP  string
}

func NewLogger(level string, logger *log.Logger, r *http.Request) (jsonlog.Logger, error) {
	requestID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	request := RequestLogEntry{
		RequestID: requestID.String(),
		RemoteIP:  r.RemoteAddr,
	}

	return leveledLogger(level, request, logger), nil
}
func leveledLogger(level string, request RequestLogEntry, logger *log.Logger) jsonlog.Logger {
	switch level {
	case "DEBUG":
		return jsonlog.NewDebugLogger(logger, request)
	case "INFO":
		return jsonlog.NewInfoLogger(logger, request)
	case "WARNING":
		return jsonlog.NewWarningLogger(logger, request)
	default:
		return jsonlog.NewErrorLogger(logger, request)
	}
}
