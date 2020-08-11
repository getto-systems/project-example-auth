package infra

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger interface {
		ParseTicketLogger

		IssueTicketLogger
		IssueApiTokenLogger
		IssueContentTokenLogger
	}

	ParseTicketLogger interface {
		TryToParseTicket(request.Request)
		FailedToParseTicket(request.Request, error)
		FailedToParseTicketBecauseNonceMatchFailed(request.Request, error)
		ParseTicket(request.Request, user.User)
	}

	IssueTicketLogger interface {
		TryToIssueTicket(request.Request, user.User, credential.TicketExpires)
		FailedToIssueTicket(request.Request, user.User, credential.TicketExpires, error)
		IssueTicket(request.Request, user.User, credential.TicketExpires)
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, credential.TokenExpires)
		FailedToIssueApiToken(request.Request, user.User, credential.TokenExpires, error)
		IssueApiToken(request.Request, user.User, credential.ApiRoles, credential.TokenExpires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, credential.TokenExpires)
		FailedToIssueContentToken(request.Request, user.User, credential.TokenExpires, error)
		IssueContentToken(request.Request, user.User, credential.TokenExpires)
	}
)
