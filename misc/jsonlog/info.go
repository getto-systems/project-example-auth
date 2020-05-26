package jsonlog

type InfoLogger struct {
	config config
}

func NewInfoLogger(logger logger, request interface{}) Logger {
	return InfoLogger{
		config: config{
			logger:  logger,
			request: request,
		},
	}
}

func (logger InfoLogger) Audit(v ...interface{}) {
	logger.config.Audit(v...)
}
func (logger InfoLogger) Auditf(format string, v ...interface{}) {
	logger.config.Auditf(format, v...)
}

func (logger InfoLogger) Error(v ...interface{}) {
	logger.config.Error(v...)
}
func (logger InfoLogger) Errorf(format string, v ...interface{}) {
	logger.config.Errorf(format, v...)
}

func (logger InfoLogger) Warning(v ...interface{}) {
	logger.config.Warning(v...)
}
func (logger InfoLogger) Warningf(format string, v ...interface{}) {
	logger.config.Warningf(format, v...)
}

func (logger InfoLogger) Info(v ...interface{}) {
	logger.config.Info(v...)
}
func (logger InfoLogger) Infof(format string, v ...interface{}) {
	logger.config.Infof(format, v...)
}

func (InfoLogger) Debug(v ...interface{}) {
	// noop
}
func (InfoLogger) Debugf(format string, v ...interface{}) {
	// noop
}
