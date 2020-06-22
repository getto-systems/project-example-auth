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
	Message           string
	Request           *data.Request
	AuthenticatedUser *data.AuthenticatedUser
	Profile           *data.Profile
	Resource          *data.Resource
	Error             error
}

func (logger UserLogger) Authenticated(request data.Request, userID data.UserID, profile data.Profile) {
	logger.logger.Audit(Log{
		Message:           "authenticated",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Profile:           &profile,
	})
}

func (logger UserLogger) Authorized(request data.Request, userID data.UserID, profile data.Profile, resource data.Resource) {
	logger.logger.Audit(Log{
		Message:           "authorized",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Profile:           &profile,
		Resource:          &resource,
	})
}

func (logger UserLogger) TicketRenewed(request data.Request, userID data.UserID) {
	logger.logger.Audit(Log{
		Message:           "ticket renewed",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
	})
}

func (logger UserLogger) PasswordMatchFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Audit(Log{
		Message:           "password match failed",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Error:             err,
	})
}

func (logger UserLogger) AuthorizeFailed(request data.Request, userID data.UserID, resource data.Resource, err error) {
	logger.logger.Audit(Log{
		Message:           "authorize failed",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Resource:          &resource,
		Error:             err,
	})
}

func (logger UserLogger) TicketRenewFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Info(Log{
		Message:           "ticket renew failed",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Error:             err,
	})
}

func (logger UserLogger) TicketIssueFailed(request data.Request, userID data.UserID, err error) {
	logger.logger.Info(Log{
		Message:           "ticket issue failed",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
		Error:             err,
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

func (logger UserLogger) TicketRenewing(request data.Request, userID data.UserID) {
	logger.logger.Debug(Log{
		Message:           "ticket renewing",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
	})
}

func (logger UserLogger) PasswordMatching(request data.Request, userID data.UserID) {
	logger.logger.Debug(Log{
		Message:           "password matching",
		Request:           &request,
		AuthenticatedUser: authenticatedUser(userID),
	})
}

func (logger UserLogger) Authorizing(request data.Request, resource data.Resource) {
	logger.logger.Debug(Log{
		Message:  "authorizing",
		Request:  &request,
		Resource: &resource,
	})
}

func authenticatedUser(userID data.UserID) *data.AuthenticatedUser {
	return &data.AuthenticatedUser{
		UserID: userID,
	}
}
