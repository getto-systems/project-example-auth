package simple_logger

import (
	"log"
	"os"
)

type WarningLogger struct {
	logger *log.Logger
}

func NewWarningLogger() WarningLogger {
	return WarningLogger{logger: log.New(os.Stdout, "", 0)}
}

func (logger WarningLogger) Audit(v ...interface{}) {
	message(logger.logger, "audit", v...)
}
func (logger WarningLogger) Auditf(format string, v ...interface{}) {
	messagef(logger.logger, "audit", format, v...)
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
	message(logger.logger, "warning", v...)
}
func (logger WarningLogger) Warningf(format string, v ...interface{}) {
	messagef(logger.logger, "warning", format, v...)
}

func (logger WarningLogger) Error(v ...interface{}) {
	message(logger.logger, "error", v...)
}
func (logger WarningLogger) Errorf(format string, v ...interface{}) {
	messagef(logger.logger, "error", format, v...)
}
