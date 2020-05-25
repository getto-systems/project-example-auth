package simple_logger

import (
	"log"
	"os"
)

type InfoLogger struct {
	logger  *log.Logger
	request interface{}
}

func (logger InfoLogger) Logger() *log.Logger {
	return logger.logger
}

func (logger InfoLogger) Request() interface{} {
	return logger.request
}

func NewInfoLogger(request interface{}) InfoLogger {
	return InfoLogger{
		logger:  log.New(os.Stdout, "", 0),
		request: request,
	}
}

func (logger InfoLogger) Audit(v ...interface{}) {
	message(logger, "audit", v...)
}
func (logger InfoLogger) Auditf(format string, v ...interface{}) {
	messagef(logger, "audit", format, v...)
}

func (InfoLogger) Debug(v ...interface{}) {
	// noop
}
func (InfoLogger) Debugf(format string, v ...interface{}) {
	// noop
}

func (logger InfoLogger) Info(v ...interface{}) {
	message(logger, "info", v...)
}
func (logger InfoLogger) Infof(format string, v ...interface{}) {
	messagef(logger, "info", format, v...)
}

func (logger InfoLogger) Warning(v ...interface{}) {
	message(logger, "warning", v...)
}
func (logger InfoLogger) Warningf(format string, v ...interface{}) {
	messagef(logger, "warning", format, v...)
}

func (logger InfoLogger) Error(v ...interface{}) {
	message(logger, "error", v...)
}
func (logger InfoLogger) Errorf(format string, v ...interface{}) {
	messagef(logger, "error", format, v...)
}
