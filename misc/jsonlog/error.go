package jsonlog

type ErrorLogger struct {
	config config
}

func NewErrorLogger(logger logger, request interface{}) Logger {
	return ErrorLogger{
		config: config{
			logger:  logger,
			request: request,
		},
	}
}

func (logger ErrorLogger) Audit(v ...interface{}) {
	logger.config.Audit(v...)
}
func (logger ErrorLogger) Auditf(format string, v ...interface{}) {
	logger.config.Auditf(format, v...)
}

func (logger ErrorLogger) Error(v ...interface{}) {
	logger.config.Error(v...)
}
func (logger ErrorLogger) Errorf(format string, v ...interface{}) {
	logger.config.Errorf(format, v...)
}

func (ErrorLogger) Warning(v ...interface{}) {
	// noop
}
func (ErrorLogger) Warningf(format string, v ...interface{}) {
	// noop
}

func (ErrorLogger) Info(v ...interface{}) {
	// noop
}
func (ErrorLogger) Infof(format string, v ...interface{}) {
	// noop
}

func (ErrorLogger) Debug(v ...interface{}) {
	// noop
}
func (ErrorLogger) Debugf(format string, v ...interface{}) {
	// noop
}
