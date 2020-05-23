package simple_logger

import (
	"log"
	"os"
)

type InfoLogger struct {
	logger *log.Logger
}

func NewInfoLogger() InfoLogger {
	return InfoLogger{logger: log.New(os.Stdout, "", 0)}
}

func (logger InfoLogger) Audit(v ...interface{}) {
	message(logger.logger, "audit", v...)
}
func (logger InfoLogger) Auditf(format string, v ...interface{}) {
	messagef(logger.logger, "audit", format, v...)
}

func (InfoLogger) Debug(v ...interface{}) {
	// noop
}
func (InfoLogger) Debugf(format string, v ...interface{}) {
	// noop
}

func (logger InfoLogger) Info(v ...interface{}) {
	message(logger.logger, "info", v...)
}
func (logger InfoLogger) Infof(format string, v ...interface{}) {
	messagef(logger.logger, "info", format, v...)
}

func (logger InfoLogger) Warning(v ...interface{}) {
	message(logger.logger, "warning", v...)
}
func (logger InfoLogger) Warningf(format string, v ...interface{}) {
	messagef(logger.logger, "warning", format, v...)
}

func (logger InfoLogger) Error(v ...interface{}) {
	message(logger.logger, "error", v...)
}
func (logger InfoLogger) Errorf(format string, v ...interface{}) {
	messagef(logger.logger, "error", format, v...)
}
