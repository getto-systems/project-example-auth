package subscriber

import (
	"github.com/getto-systems/project-example-id/data"
)

type Logger interface {
	Audit(Log)
	Info(Log)
	Debug(Log)
}

type UserLogger struct {
	logger Logger
}

func NewUserLogger(logger Logger) UserLogger {
	return UserLogger{
		logger: logger,
	}
}

type Log struct {
	Message  string
	Request  *data.Request
	User     *data.User
	Ticket   *data.Ticket
	Resource *data.Resource
	Error    error
}

func (logger UserLogger) Authenticated(request data.Request, ticket data.Ticket) {
	logger.logger.Audit(Log{
		Message: "authenticated",
		Request: &request,
		Ticket:  &ticket,
	})
}

func (logger UserLogger) Authorized(request data.Request, ticket data.Ticket, resource data.Resource) {
	logger.logger.Audit(Log{
		Message:  "authorized",
		Request:  &request,
		Ticket:   &ticket,
		Resource: &resource,
	})
}

func (logger UserLogger) AuthorizeFailed(request data.Request, ticket data.Ticket, resource data.Resource, err error) {
	logger.logger.Audit(Log{
		Message:  "authorize failed",
		Request:  &request,
		Ticket:   &ticket,
		Resource: &resource,
		Error:    err,
	})
}

func (logger UserLogger) TicketRenewed(request data.Request, ticket data.Ticket) {
	logger.logger.Audit(Log{
		Message: "ticket renewed",
		Request: &request,
		Ticket:  &ticket,
	})
}

func (logger UserLogger) PasswordMatchFailed(request data.Request, user data.User, err error) {
	logger.logger.Audit(Log{
		Message: "password match failed",
		Request: &request,
		User:    &user,
		Error:   err,
	})
}

func (logger UserLogger) TicketRenewFailed(request data.Request, ticket data.Ticket, err error) {
	logger.logger.Info(Log{
		Message: "ticket renew failed",
		Request: &request,
		Ticket:  &ticket,
		Error:   err,
	})
}

func (logger UserLogger) TicketIssueFailed(request data.Request, user data.User, err error) {
	logger.logger.Info(Log{
		Message: "ticket issue failed",
		Request: &request,
		User:    &user,
		Error:   err,
	})
}

func (logger UserLogger) TicketRenewing(request data.Request, ticket data.Ticket) {
	logger.logger.Debug(Log{
		Message: "ticket renewing",
		Request: &request,
		Ticket:  &ticket,
	})
}

func (logger UserLogger) PasswordMatching(request data.Request, user data.User) {
	logger.logger.Debug(Log{
		Message: "password matching",
		Request: &request,
		User:    &user,
	})
}

func (logger UserLogger) Authorizing(request data.Request, resource data.Resource) {
	logger.logger.Debug(Log{
		Message:  "authorizing",
		Request:  &request,
		Resource: &resource,
	})
}

func (logger UserLogger) AuthorizeTokenParseFailed(request data.Request, resource data.Resource, err error) {
	logger.logger.Debug(Log{
		Message:  "authorize token parse failed",
		Request:  &request,
		Resource: &resource,
		Error:    err,
	})
}
