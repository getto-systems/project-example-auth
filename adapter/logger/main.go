package logger

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"

	"github.com/getto-systems/applog-go"
)

type LogEntry struct {
	Time      string `json:"time"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	RemoteIP  string `json:"remote_ip"`
	RequestID string `json:"request_id"`
}

func NewLogger(level string, logger *log.Logger, r *http.Request) (applog.Logger, error) {
	randomID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	requestID := randomID.String()
	remoteIP := r.RemoteAddr

	entry := func(level string, message string) string {
		data, err := json.Marshal(LogEntry{
			Time:      time.Now().UTC().String(),
			Level:     level,
			Message:   message,
			RemoteIP:  remoteIP,
			RequestID: requestID,
		})
		if err != nil {
			return err.Error()
		}

		return string(data)
	}

	return leveledLogger(level, logger, entry), nil
}
func leveledLogger(level string, output applog.Output, entry applog.Entry) applog.Logger {
	switch level {
	case "DEBUG":
		return applog.NewDebugLogger(output, entry)
	case "INFO":
		return applog.NewInfoLogger(output, entry)
	case "WARNING":
		return applog.NewWarnLogger(output, entry)
	default:
		return applog.NewErrorLogger(output, entry)
	}
}
