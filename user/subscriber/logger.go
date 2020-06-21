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
	Request  data.Request
	UserID   data.UserID
	Resource data.Resource
	Error    error
}

func (logger UserLogger) Authenticated(request data.Request, userID data.UserID) {
	logger.logger.Audit(Log{
		Message: "authenticated",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) Authorized(request data.Request, userID data.UserID, resource data.Resource) {
	logger.logger.Audit(Log{
		Message:  "authorized",
		Request:  request,
		UserID:   userID,
		Resource: resource,
	})
}

func (logger UserLogger) TicketRenewed(request data.Request, userID data.UserID) {
	logger.logger.Audit(Log{
		Message: "ticket renewed",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) PasswordMatchFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Audit(Log{
		Message: "password match failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) AuthorizeFailed(request data.Request, userID data.UserID, resource data.Resource, err error) {
	logger.logger.Audit(Log{
		Message:  "authorize failed",
		Request:  request,
		UserID:   userID,
		Resource: resource,
		Error:    err,
	})
}

func (logger UserLogger) TicketRenewFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Info(Log{
		Message: "ticket renew failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) TicketIssueFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Info(Log{
		Message: "ticket issue failed",
		Request: request,
		UserID:  userID,
		Error:   err,
	})
}

func (logger UserLogger) AuthorizeTokenParseFailed(request data.Request, resource data.Resource, err error) {
	logger.logger.Debug(Log{
		Message:  "authorize token parse failed",
		Request:  request,
		Resource: resource,
		Error:    err,
	})
}

func (logger UserLogger) TicketRenewing(request data.Request, userID data.UserID) {
	logger.logger.Debug(Log{
		Message: "ticket renewing",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) PasswordMatching(request data.Request, userID data.UserID) {
	logger.logger.Debug(Log{
		Message: "password matching",
		Request: request,
		UserID:  userID,
	})
}

func (logger UserLogger) Authorizing(request data.Request, resource data.Resource) {
	logger.logger.Debug(Log{
		Message:  "authorizing",
		Request:  request,
		Resource: resource,
	})
}
