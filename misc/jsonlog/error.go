package jsonlog

type ErrorLogger struct {
	errorConfig config
}

func (logger ErrorLogger) config() config {
	return logger.errorConfig
}

func NewErrorLogger(logger logger, request interface{}) ErrorLogger {
	return ErrorLogger{
		errorConfig: config{
			logger:  logger,
			request: request,
		},
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
