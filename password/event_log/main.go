package event_log

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

type EventLogger struct {
	logger event_log.Logger
}

func NewEventLogger(logger event_log.Logger) EventLogger {
	return EventLogger{
		logger: logger,
	}
}

func (log EventLogger) handler() password.EventHandler {
	return log
}

func (log EventLogger) ValidatePassword(request data.Request, login password.Login) {
	log.logger.Debug(event_log.Entry{
		Message: "validate password",
		Request: request,
		Login:   &login,
	})
}

func (log EventLogger) ValidatePasswordFailed(request data.Request, login password.Login, err error) {
	log.logger.Audit(event_log.Entry{
		Message: "validate password failed",
		Request: request,
		Login:   &login,
		Error:   err,
	})
}

func (log EventLogger) AuthenticatedByPassword(request data.Request, login password.Login, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "authenticated by password",
		Request: request,
		Login:   &login,
		User:    &user,
	})
}

func (log EventLogger) GetLogin(request data.Request, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "get login",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) GetLoginFailed(request data.Request, user data.User, err error) {
	log.logger.Info(event_log.Entry{
		Message: "get login failed",
		Request: request,
		User:    &user,
		Error:   err,
	})
}

func (log EventLogger) RegisterPassword(request data.Request, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "register password",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) RegisterPasswordFailed(request data.Request, user data.User, err error) {
	log.logger.Info(event_log.Entry{
		Message: "register password failed",
		Request: request,
		User:    &user,
		Error:   err,
	})
}

func (log EventLogger) RegisteredPassword(request data.Request, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "password registered",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) IssueReset(request data.Request, login password.Login, expires data.Expires) {
	log.logger.Debug(event_log.Entry{
		Message: "issue reset",
		Request: request,
		Login:   &login,
		Expires: &expires,
	})
}

func (log EventLogger) IssueResetFailed(request data.Request, login password.Login, expires data.Expires, err error) {
	log.logger.Info(event_log.Entry{
		Message: "issue reset failed",
		Request: request,
		Login:   &login,
		Expires: &expires,
		Error:   err,
	})
}

func (log EventLogger) IssuedReset(request data.Request, login password.Login, expires data.Expires, reset password.Reset, _ password.ResetToken) {
	log.logger.Debug(event_log.Entry{
		Message: "issue reset",
		Request: request,
		Reset:   &reset,
		Login:   &login,
		Expires: &expires,
	})
}

func (log EventLogger) GetResetStatus(request data.Request, reset password.Reset) {
	log.logger.Debug(event_log.Entry{
		Message: "get reset status",
		Request: request,
		Reset:   &reset,
	})
}

func (log EventLogger) GetResetStatusFailed(request data.Request, reset password.Reset, err error) {
	log.logger.Info(event_log.Entry{
		Message: "get reset status failed",
		Request: request,
		Reset:   &reset,
		Error:   err,
	})
}

func (log EventLogger) ValidateResetToken(request data.Request) {
	log.logger.Debug(event_log.Entry{
		Message: "validate reset token",
		Request: request,
	})
}

func (log EventLogger) ValidateResetTokenFailed(request data.Request, err error) {
	log.logger.Audit(event_log.Entry{
		Message: "validate reset token failed",
		Request: request,
	})
}

func (log EventLogger) AuthenticatedByResetToken(request data.Request, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "authenticated by reset token",
		Request: request,
		User:    &user,
	})
}
