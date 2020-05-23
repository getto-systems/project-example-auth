package simple_logger

import (
	"log"
	"os"
)

type ErrorLogger struct {
	logger *log.Logger
}

func NewErrorLogger() ErrorLogger {
	return ErrorLogger{logger: log.New(os.Stdout, "", 0)}
}

func (logger ErrorLogger) Audit(v ...interface{}) {
	message(logger.logger, "audit", v...)
}
func (logger ErrorLogger) Auditf(format string, v ...interface{}) {
	messagef(logger.logger, "audit", format, v...)
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
	message(logger.logger, "error", v...)
}
func (logger ErrorLogger) Errorf(format string, v ...interface{}) {
	messagef(logger.logger, "error", format, v...)
}
