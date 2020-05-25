package simple_logger

import (
	"log"
	"os"
)

type WarningLogger struct {
	logger  *log.Logger
	request interface{}
}

func (logger WarningLogger) Logger() *log.Logger {
	return logger.logger
}

func (logger WarningLogger) Request() interface{} {
	return logger.request
}

func NewWarningLogger(request interface{}) WarningLogger {
	return WarningLogger{
		logger:  log.New(os.Stdout, "", 0),
		request: request,
	}
}

func (logger WarningLogger) Audit(v ...interface{}) {
	message(logger, "audit", v...)
}
func (logger WarningLogger) Auditf(format string, v ...interface{}) {
	messagef(logger, "audit", format, v...)
}

func (WarningLogger) Debug(v ...interface{}) {
	// noop
}
func (WarningLogger) Debugf(format string, v ...interface{}) {
	// noop
}

func (WarningLogger) Info(v ...interface{}) {
	// noop
}
func (WarningLogger) Infof(format string, v ...interface{}) {
	// noop
}

func (logger WarningLogger) Warning(v ...interface{}) {
	message(logger, "warning", v...)
}
func (logger WarningLogger) Warningf(format string, v ...interface{}) {
	messagef(logger, "warning", format, v...)
}

func (logger WarningLogger) Error(v ...interface{}) {
	message(logger, "error", v...)
}
func (logger WarningLogger) Errorf(format string, v ...interface{}) {
	messagef(logger, "error", format, v...)
}
