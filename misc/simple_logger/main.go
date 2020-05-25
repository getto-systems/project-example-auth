package simple_logger

import (
	"encoding/json"
	"log"

	"fmt"
	"time"
)

type LogEntry struct {
	Time    string      `json:"time"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Request interface{} `json:"request"`
}

type Logger interface {
	Logger() *log.Logger
	Request() interface{}
}

func message(logger Logger, level string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprint(v...))
}
func messagef(logger Logger, level string, format string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprintf(format, v...))
}
func jsonMessage(logger Logger, level string, message string) {
	data, err := json.Marshal(LogEntry{
		Time:    time.Now().UTC().String(),
		Level:   level,
		Message: message,
		Request: logger.Request(),
	})
	if err != nil {
		logger.Logger().Print(err)
		return
	}

	logger.Logger().Print(string(data))
}
