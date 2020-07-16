package event_log

import (
	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/data"
)

type PasswordEventLogger struct {
	logger Logger
}

func NewPasswordEventLogger(logger Logger) PasswordEventLogger {
	return PasswordEventLogger{
		logger: logger,
	}
}

func (log PasswordEventLogger) Handler() password.EventHandler {
	return log
}

func (log PasswordEventLogger) ValidatePassword(request data.Request, user data.User) {
	log.logger.Debug(Entry{
		Message: "validate password",
		Request: request,
		User:    &user,
	})
}

func (log PasswordEventLogger) ValidatePasswordFailed(request data.Request, user data.User, err error) {
	log.logger.Info(Entry{
		Message: "validate password failed",
		Request: request,
		User:    &user,
		Error:   err,
	})
}

func (log PasswordEventLogger) PasswordRegistered(request data.Request, user data.User) {
	log.logger.Audit(Entry{
		Message: "password registered",
		Request: request,
		User:    &user,
	})
}

func (log PasswordEventLogger) VerifyPassword(request data.Request, user data.User) {
	log.logger.Debug(Entry{
		Message: "verify password",
		Request: request,
		User:    &user,
	})
}

func (log PasswordEventLogger) VerifyPasswordFailed(request data.Request, user data.User, err error) {
	log.logger.Audit(Entry{
		Message: "verify password failed",
		Request: request,
		User:    &user,
		Error:   err,
	})
}

func (log PasswordEventLogger) AuthenticatedByPassword(request data.Request, user data.User) {
	log.logger.Audit(Entry{
		Message: "authenticated by password",
		Request: request,
		User:    &user,
	})
}
