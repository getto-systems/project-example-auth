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

func (log EventLogger) GetLogin(request data.Request, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "get login",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) LoginNotFound(request data.Request, user data.User, err error) {
	log.logger.Info(event_log.Entry{
		Message: "login not found",
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

func (log EventLogger) PasswordRegistered(request data.Request, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "password registered",
		Request: request,
		User:    &user,
	})
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
