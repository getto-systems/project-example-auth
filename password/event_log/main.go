package event_log

import (
	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
)

type EventLogger struct {
	logger event_log.Logger
}

func NewEventLogger(logger event_log.Logger) EventLogger {
	return EventLogger{
		logger: logger,
	}
}

func (log EventLogger) Handler() password.EventHandler {
	return log
}

func (log EventLogger) RegisterPassword(request data.Request, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "validate password",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) RegisterPasswordFailed(request data.Request, user data.User, err error) {
	log.logger.Info(event_log.Entry{
		Message: "validate password failed",
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

func (log EventLogger) VerifyPassword(request data.Request, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "verify password",
		Request: request,
		User:    &user,
	})
}

func (log EventLogger) VerifyPasswordFailed(request data.Request, user data.User, err error) {
	log.logger.Audit(event_log.Entry{
		Message: "verify password failed",
		Request: request,
		User:    &user,
		Error:   err,
	})
}

func (log EventLogger) AuthenticatedByPassword(request data.Request, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "authenticated by password",
		Request: request,
		User:    &user,
	})
}
