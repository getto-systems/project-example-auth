package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/getto-systems/applog-go/v2"
)

type (
	LeveledLogger struct {
		logger applog.Logger
	}

	LogEntry struct {
		Level string      `json:"level"`
		Time  string      `json:"time"`
		Entry interface{} `json:"entry"`
	}
)

func NewLeveledLogger(level string) LeveledLogger {
	return LeveledLogger{
		logger: leveledLogger(level),
	}
}
func leveledLogger(level string) applog.Logger {
	logger := log.New(os.Stdout, "", 0)

	switch level {
	case "QUIET":
		return applog.NewQuietLogger(logger)
	case "INFO":
		return applog.NewInfoLogger(logger)
	default:
		return applog.NewDebugLogger(logger)
	}
}

func (logger LeveledLogger) Audit(entry interface{}) {
	logger.logger.Always(jsonLogEntry("AUDIT", entry))
}
func (logger LeveledLogger) Error(entry interface{}) {
	logger.logger.Always(jsonLogEntry("ERROR", entry))
}
func (logger LeveledLogger) Info(entry interface{}) {
	logger.logger.Info(jsonLogEntry("INFO", entry))
}
func (logger LeveledLogger) Debug(entry interface{}) {
	logger.logger.Debug(jsonLogEntry("DEBUG", entry))
}

func jsonLogEntry(level string, entry interface{}) string {
	log := LogEntry{
		Time:  time.Now().UTC().String(),
		Level: level,
		Entry: entry,
	}

	data, err := json.Marshal(log)
	if err != nil {
		return "json marshal error"
	}

	return string(data)
}
