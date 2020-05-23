package simple_logger

import (
	"log"
	"os"
)

type DebugLogger struct {
	logger *log.Logger
}

func NewDebugLogger() DebugLogger {
	return DebugLogger{logger: log.New(os.Stdout, "", 0)}
}

func (logger DebugLogger) Audit(v ...interface{}) {
	message(logger.logger, "audit", v...)
}
func (logger DebugLogger) Auditf(format string, v ...interface{}) {
	messagef(logger.logger, "audit", format, v...)
}

func (logger DebugLogger) Debug(v ...interface{}) {
	message(logger.logger, "debug", v...)
}
func (logger DebugLogger) Debugf(format string, v ...interface{}) {
	messagef(logger.logger, "debug", format, v...)
}

func (logger DebugLogger) Info(v ...interface{}) {
	message(logger.logger, "info", v...)
}
func (logger DebugLogger) Infof(format string, v ...interface{}) {
	messagef(logger.logger, "info", format, v...)
}

func (logger DebugLogger) Warning(v ...interface{}) {
	message(logger.logger, "warning", v...)
}
func (logger DebugLogger) Warningf(format string, v ...interface{}) {
	messagef(logger.logger, "warning", format, v...)
}

func (logger DebugLogger) Error(v ...interface{}) {
	message(logger.logger, "error", v...)
}
func (logger DebugLogger) Errorf(format string, v ...interface{}) {
	messagef(logger.logger, "error", format, v...)
}
