package simple_logger

import (
	"log"
	"os"
)

type ErrorLogger struct {
	logger  *log.Logger
	request interface{}
}

func (logger ErrorLogger) Logger() *log.Logger {
	return logger.logger
}

func (logger ErrorLogger) Request() interface{} {
	return logger.request
}

func NewErrorLogger(request interface{}) ErrorLogger {
	return ErrorLogger{
		logger:  log.New(os.Stdout, "", 0),
		request: request,
	}
}

func (logger ErrorLogger) Audit(v ...interface{}) {
	message(logger, "audit", v...)
}
func (logger ErrorLogger) Auditf(format string, v ...interface{}) {
	messagef(logger, "audit", format, v...)
}

func (ErrorLogger) Debug(v ...interface{}) {
	// noop
}
func (ErrorLogger) Debugf(format string, v ...interface{}) {
	// noop
}

func (ErrorLogger) Info(v ...interface{}) {
	// noop
}
func (ErrorLogger) Infof(format string, v ...interface{}) {
	// noop
}

func (ErrorLogger) Warning(v ...interface{}) {
	// noop
}
func (ErrorLogger) Warningf(format string, v ...interface{}) {
	// noop
}

func (logger ErrorLogger) Error(v ...interface{}) {
	message(logger, "error", v...)
}
func (logger ErrorLogger) Errorf(format string, v ...interface{}) {
	messagef(logger, "error", format, v...)
}
