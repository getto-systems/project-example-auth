package event_log

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/ticket"
)

type EventLogger struct {
	logger event_log.Logger
}

func NewEventLogger(logger event_log.Logger) EventLogger {
	return EventLogger{
		logger: logger,
	}
}

func (log EventLogger) handler() ticket.EventHandler {
	return log
}

func (log EventLogger) IssueApiToken(request data.Request, user data.User, roles data.Roles, expires data.Expires) {
	log.logger.Debug(event_log.Entry{
		Message: "issue api token",
		Request: request,
		User:    &user,
		Roles:   &roles,
		Expires: &expires,
	})
}

func (log EventLogger) IssueApiTokenFailed(request data.Request, user data.User, roles data.Roles, expires data.Expires, err error) {
	log.logger.Info(event_log.Entry{
		Message: "issue api token failed",
		Request: request,
		User:    &user,
		Roles:   &roles,
		Expires: &expires,
		Error:   err,
	})
}

func (log EventLogger) IssueContentToken(request data.Request, user data.User, expires data.Expires) {
	log.logger.Debug(event_log.Entry{
		Message: "issue content token",
		Request: request,
		User:    &user,
		Expires: &expires,
	})
}

func (log EventLogger) IssueContentTokenFailed(request data.Request, user data.User, expires data.Expires, err error) {
	log.logger.Info(event_log.Entry{
		Message: "issue content token failed",
		Request: request,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log EventLogger) ExtendTicket(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires) {
	log.logger.Debug(event_log.Entry{
		Message: "extend ticket",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Expires: &expires,
	})
}

func (log EventLogger) ExtendTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires, err error) {
	log.logger.Info(event_log.Entry{
		Message: "extend ticket failed",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log EventLogger) IssueTicket(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit) {
	log.logger.Debug(event_log.Entry{
		Message: "issue ticket",
		Request: request,
		User:    &user,
		Expires: &expires,
	})
}

func (log EventLogger) IssueTicketFailed(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit, err error) {
	log.logger.Info(event_log.Entry{
		Message: "issue ticket failed",
		Request: request,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log EventLogger) ShrinkTicket(request data.Request, nonce ticket.Nonce, user data.User) {
	log.logger.Debug(event_log.Entry{
		Message: "shrink ticket",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
	})
}

func (log EventLogger) ShrinkTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, err error) {
	log.logger.Info(event_log.Entry{
		Message: "shrink ticket failed",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Error:   err,
	})
}

func (log EventLogger) ValidateTicket(request data.Request) {
	log.logger.Debug(event_log.Entry{
		Message: "validate ticket",
		Request: request,
	})
}

func (log EventLogger) ValidateTicketFailed(request data.Request, err error) {
	log.logger.Info(event_log.Entry{
		Message: "validate ticket failed",
		Request: request,
		Error:   err,
	})
}

func (log EventLogger) AuthenticatedByTicket(request data.Request, user data.User) {
	log.logger.Audit(event_log.Entry{
		Message: "authenticated by ticket",
		Request: request,
		User:    &user,
	})
}
