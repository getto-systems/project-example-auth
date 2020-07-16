package event_log

import (
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type TicketEventLogger struct {
	logger Logger
}

func NewTicketEventLogger(logger Logger) TicketEventLogger {
	return TicketEventLogger{
		logger: logger,
	}
}

func (log TicketEventLogger) Handler() ticket.EventHandler {
	return log
}

func (log TicketEventLogger) IssueApiToken(request data.Request, user data.User, roles data.Roles, expires data.Expires) {
	log.logger.Debug(Entry{
		Message: "issue api token",
		Request: request,
		User:    &user,
		Roles:   &roles,
		Expires: &expires,
	})
}

func (log TicketEventLogger) IssueApiTokenFailed(request data.Request, user data.User, roles data.Roles, expires data.Expires, err error) {
	log.logger.Info(Entry{
		Message: "issue api token failed",
		Request: request,
		User:    &user,
		Roles:   &roles,
		Expires: &expires,
		Error:   err,
	})
}

func (log TicketEventLogger) IssueContentToken(request data.Request, user data.User, expires data.Expires) {
	log.logger.Debug(Entry{
		Message: "issue content token",
		Request: request,
		User:    &user,
		Expires: &expires,
	})
}

func (log TicketEventLogger) IssueContentTokenFailed(request data.Request, user data.User, expires data.Expires, err error) {
	log.logger.Info(Entry{
		Message: "issue content token failed",
		Request: request,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log TicketEventLogger) ExtendTicket(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires) {
	log.logger.Debug(Entry{
		Message: "extend ticket",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Expires: &expires,
	})
}

func (log TicketEventLogger) ExtendTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires, err error) {
	log.logger.Info(Entry{
		Message: "extend ticket failed",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log TicketEventLogger) IssueTicket(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit) {
	log.logger.Debug(Entry{
		Message: "issue ticket",
		Request: request,
		User:    &user,
		Expires: &expires,
	})
}

func (log TicketEventLogger) IssueTicketFailed(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit, err error) {
	log.logger.Info(Entry{
		Message: "issue ticket failed",
		Request: request,
		User:    &user,
		Expires: &expires,
		Error:   err,
	})
}

func (log TicketEventLogger) ShrinkTicket(request data.Request, nonce ticket.Nonce, user data.User) {
	log.logger.Debug(Entry{
		Message: "shrink ticket",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
	})
}

func (log TicketEventLogger) ShrinkTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, err error) {
	log.logger.Info(Entry{
		Message: "shrink ticket failed",
		Request: request,
		Nonce:   &nonce,
		User:    &user,
		Error:   err,
	})
}

func (log TicketEventLogger) VerifyTicket(request data.Request) {
	log.logger.Debug(Entry{
		Message: "verify ticket",
		Request: request,
	})
}

func (log TicketEventLogger) VerifyTicketFailed(request data.Request, err error) {
	log.logger.Info(Entry{
		Message: "verify ticket failed",
		Request: request,
		Error:   err,
	})
}

func (log TicketEventLogger) AuthenticatedByTicket(request data.Request, user data.User) {
	log.logger.Audit(Entry{
		Message: "authenticated by ticket",
		Request: request,
		User:    &user,
	})
}
