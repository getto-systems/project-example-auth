package subscriber

import (
	"github.com/getto-systems/project-example-id/basic"
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
	Request  basic.Request
	UserID   basic.UserID
	Resource basic.Resource
	Error    error
}

func (logger UserLogger) Authenticated(request basic.Request, userID basic.UserID) {
	logger.logger.Audit(Log{
		Message: "authenticated",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) Authorized(request basic.Request, userID basic.UserID, resource basic.Resource) {
	logger.logger.Audit(Log{
		Message:  "authorized",
		Request:  request,
		UserID:   userID,
		Resource: resource,
	})
}

func (logger UserLogger) TicketRenewed(request basic.Request, userID basic.UserID) {
	logger.logger.Audit(Log{
		Message: "ticket renewed",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) PasswordMatchFailed(request basic.Request, userID basic.UserID, err error) {
	logger.logger.Audit(Log{
		Message: "password match failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) AuthorizeFailed(request basic.Request, userID basic.UserID, resource basic.Resource, err error) {
	logger.logger.Audit(Log{
		Message:  "authorize failed",
		Request:  request,
		UserID:   userID,
		Resource: resource,
		Error:    err,
	})
}

func (logger UserLogger) TicketRenewFailed(request basic.Request, userID basic.UserID, err error) {
	logger.logger.Info(Log{
		Message: "ticket renew failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) TicketIssueFailed(request basic.Request, userID basic.UserID, err error) {
	logger.logger.Info(Log{
		Message: "ticket issue failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) AuthorizeTokenParseFailed(request basic.Request, resource basic.Resource, err error) {
	logger.logger.Debug(Log{
		Message:  "authorize token parse failed",
		Request:  request,
		Resource: resource,
		Error:    err,
	})
}

func (logger UserLogger) TicketRenewing(request basic.Request, userID basic.UserID) {
	logger.logger.Debug(Log{
		Message: "ticket renewing",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) PasswordMatching(request basic.Request, userID basic.UserID) {
	logger.logger.Debug(Log{
		Message: "password matching",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) Authorizing(request basic.Request, resource basic.Resource) {
	logger.logger.Debug(Log{
		Message:  "authorizing",
		Request:  request,
		Resource: resource,
	})
}
