package simple_logger

import (
	"encoding/json"
	"log"

	"fmt"
	"time"
)

type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func message(logger *log.Logger, level string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprint(v...))
}
func messagef(logger *log.Logger, level string, format string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprintf(format, v...))
}
func jsonMessage(logger *log.Logger, level string, message string) {
	data, err := json.Marshal(LogEntry{
		Time:    time.Now().UTC().String(),
		Level:   level,
		Message: message,
	})
	if err != nil {
		logger.Print(err)
		return
	}

	logger.Print(string(data))
}
