package logger

import (
	"github.com/google/uuid"
	"log"
	"net/http"

	"github.com/getto-systems/project-example-id/misc/simple_logger"

	"github.com/getto-systems/project-example-id/journal"
)

type RequestLogEntry struct {
	RequestID string
	RemoteIP  string
}

func NewLogger(level string, logger *log.Logger, r *http.Request) (journal.Logger, error) {
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
func leveledLogger(level string, request RequestLogEntry, logger *log.Logger) journal.Logger {
	switch level {
	case "DEBUG":
		return simple_logger.NewDebugLogger(logger, request)
	case "INFO":
		return simple_logger.NewInfoLogger(logger, request)
	case "WARNING":
		return simple_logger.NewWarningLogger(logger, request)
	default:
		return simple_logger.NewErrorLogger(logger, request)
	}
}
