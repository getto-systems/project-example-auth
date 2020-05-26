package jsonlog

type DebugLogger struct {
	debugConfig config
}

func (logger DebugLogger) config() config {
	return logger.debugConfig
}

func NewDebugLogger(logger logger, request interface{}) DebugLogger {
	return DebugLogger{
		debugConfig: config{
			logger:  logger,
			request: request,
		},
	}
}

func (logger DebugLogger) Audit(v ...interface{}) {
	message(logger, "audit", v...)
}
func (logger DebugLogger) Auditf(format string, v ...interface{}) {
	messagef(logger, "audit", format, v...)
}

func (logger DebugLogger) Debug(v ...interface{}) {
	message(logger, "debug", v...)
}
func (logger DebugLogger) Debugf(format string, v ...interface{}) {
	messagef(logger, "debug", format, v...)
}

func (logger DebugLogger) Info(v ...interface{}) {
	message(logger, "info", v...)
}
func (logger DebugLogger) Infof(format string, v ...interface{}) {
	messagef(logger, "info", format, v...)
}

func (logger DebugLogger) Warning(v ...interface{}) {
	message(logger, "warning", v...)
}
func (logger DebugLogger) Warningf(format string, v ...interface{}) {
	messagef(logger, "warning", format, v...)
}

func (logger DebugLogger) Error(v ...interface{}) {
	message(logger, "error", v...)
}
func (logger DebugLogger) Errorf(format string, v ...interface{}) {
	messagef(logger, "error", format, v...)
}
