package jsonlog

import (
	"encoding/json"

	"fmt"
	"time"
)

const (
	level_audit   = "audit"
	level_error   = "error"
	level_warning = "warning"
	level_info    = "info"
	level_debug   = "debug"
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
	Audit(v ...interface{})
	Auditf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

func (config config) Audit(v ...interface{}) {
	config.message(level_audit, v...)
}
func (config config) Auditf(format string, v ...interface{}) {
	config.messagef(level_audit, format, v...)
}

func (config config) Debug(v ...interface{}) {
	config.message(level_debug, v...)
}
func (config config) Debugf(format string, v ...interface{}) {
	config.messagef(level_debug, format, v...)
}

func (config config) Info(v ...interface{}) {
	config.message(level_info, v...)
}
func (config config) Infof(format string, v ...interface{}) {
	config.messagef(level_info, format, v...)
}

func (config config) Warning(v ...interface{}) {
	config.message(level_warning, v...)
}
func (config config) Warningf(format string, v ...interface{}) {
	config.messagef(level_warning, format, v...)
}

func (config config) Error(v ...interface{}) {
	config.message(level_error, v...)
}
func (config config) Errorf(format string, v ...interface{}) {
	config.messagef(level_error, format, v...)
}

func (config config) message(level string, v ...interface{}) {
	config.logger.Println(config.jsonMessage(level, fmt.Sprint(v...)))
}
func (config config) messagef(level string, format string, v ...interface{}) {
	config.logger.Println(config.jsonMessage(level, fmt.Sprintf(format, v...)))
}
func (config config) jsonMessage(level string, message string) string {
	data, err := json.Marshal(LogEntry{
		Time:    time.Now().UTC().String(),
		Level:   level,
		Message: message,
		Request: config.request,
	})
	if err != nil {
		return err.Error()
	}

	return string(data)
}
