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
	Message string
	Request data.Request
	User    *data.User
	Ticket  *data.Ticket
	Error   error
}

func (logger UserLogger) Authenticated(request data.Request, ticket data.Ticket) {
	logger.logger.Audit(Log{
		Message: "authenticated",
		Request: request,
		Ticket:  &ticket,
	})
}

func (logger UserLogger) TicketIssueFailed(request data.Request, ticket data.Ticket, err error) {
	logger.logger.Info(Log{
		Message: "ticket issue failed",
		Request: request,
		Ticket:  &ticket,
		Error:   err,
	})
}

func (logger UserLogger) SignedTicketParsing(request data.Request) {
	logger.logger.Debug(Log{
		Message: "signed ticket parsing",
		Request: request,
	})
}

func (logger UserLogger) SignedTicketParseFailed(request data.Request, err error) {
	logger.logger.Debug(Log{
		Message: "signed ticket parse failed",
		Request: request,
		Error:   err,
	})
}

func (logger UserLogger) PasswordMatching(request data.Request, userID data.UserID) {
	logger.logger.Debug(Log{
		Message: "password matching",
		Request: request,
		User:    &data.User{UserID: userID},
	})
}

func (logger UserLogger) PasswordMatchFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Audit(Log{
		Message: "password match failed",
		Request: request,
		User:    &data.User{UserID: userID},
		Error:   err,
	})
}
