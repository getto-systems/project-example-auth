package jsonlog

type WarningLogger struct {
	config config
}

func NewWarningLogger(logger logger, request interface{}) Logger {
	return WarningLogger{
		config: config{
			logger:  logger,
			request: request,
		},
	}
}

func (logger WarningLogger) Audit(v ...interface{}) {
	logger.config.Audit(v...)
}
func (logger WarningLogger) Auditf(format string, v ...interface{}) {
	logger.config.Auditf(format, v...)
}

func (logger WarningLogger) Error(v ...interface{}) {
	logger.config.Error(v...)
}
func (logger WarningLogger) Errorf(format string, v ...interface{}) {
	logger.config.Errorf(format, v...)
}

func (logger WarningLogger) Warning(v ...interface{}) {
	logger.config.Warning(v...)
}
func (logger WarningLogger) Warningf(format string, v ...interface{}) {
	logger.config.Warningf(format, v...)
}

func (WarningLogger) Info(v ...interface{}) {
	// noop
}
func (WarningLogger) Infof(format string, v ...interface{}) {
	// noop
}

func (WarningLogger) Debug(v ...interface{}) {
	// noop
}
func (WarningLogger) Debugf(format string, v ...interface{}) {
	// noop
}
