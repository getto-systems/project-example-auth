package jsonlog

import (
	"encoding/json"

	"fmt"
	"time"
)

type logger interface {
	Println(v ...interface{})
}

type config struct {
	logger  logger
	request interface{}
}

type LogEntry struct {
	Time    string      `json:"time"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Request interface{} `json:"request"`
}

type Logger interface {
	config() config

	Audit(v ...interface{})
	Auditf(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

func message(logger Logger, level string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprint(v...))
}
func messagef(logger Logger, level string, format string, v ...interface{}) {
	jsonMessage(logger, level, fmt.Sprintf(format, v...))
}
func jsonMessage(logger Logger, level string, message string) {
	config := logger.config()

	data, err := json.Marshal(LogEntry{
		Time:    time.Now().UTC().String(),
		Level:   level,
		Message: message,
		Request: config.request,
	})
	if err != nil {
		config.logger.Println(err)
		return
	}

	config.logger.Println(string(data))
}
