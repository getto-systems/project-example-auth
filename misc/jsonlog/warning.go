package jsonlog

type WarningLogger struct {
	warningConfig config
}

func (logger WarningLogger) config() config {
	return logger.warningConfig
}

func NewWarningLogger(logger logger, request interface{}) WarningLogger {
	return WarningLogger{
		warningConfig: config{
			logger:  logger,
			request: request,
		},
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
