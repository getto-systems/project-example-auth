package jsonlog

type DebugLogger struct {
	config config
}

func NewDebugLogger(logger logger, request interface{}) Logger {
	return DebugLogger{
		config: config{
			logger:  logger,
			request: request,
		},
	}
}

func (logger DebugLogger) Audit(v ...interface{}) {
	logger.config.Audit(v...)
}
func (logger DebugLogger) Auditf(format string, v ...interface{}) {
	logger.config.Auditf(format, v...)
}

func (logger DebugLogger) Error(v ...interface{}) {
	logger.config.Error(v...)
}
func (logger DebugLogger) Errorf(format string, v ...interface{}) {
	logger.config.Errorf(format, v...)
}

func (logger DebugLogger) Warning(v ...interface{}) {
	logger.config.Warning(v...)
}
func (logger DebugLogger) Warningf(format string, v ...interface{}) {
	logger.config.Warningf(format, v...)
}

func (logger DebugLogger) Info(v ...interface{}) {
	logger.config.Info(v...)
}
func (logger DebugLogger) Infof(format string, v ...interface{}) {
	logger.config.Infof(format, v...)
}

func (logger DebugLogger) Debug(v ...interface{}) {
	logger.config.Debug(v...)
}
func (logger DebugLogger) Debugf(format string, v ...interface{}) {
	logger.config.Debugf(format, v...)
}
