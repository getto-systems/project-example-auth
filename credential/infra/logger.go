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
		TryToParseTicket(request.Request, credential.TicketNonce)
		FailedToParseTicket(request.Request, credential.TicketNonce, error)
		FailedToParseTicketBecauseNonceMatchFailed(request.Request, credential.TicketNonce, error)
		ParseTicket(request.Request, credential.TicketNonce, user.User)
	}

	IssueTicketLogger interface {
		TryToIssueTicket(request.Request, user.User, credential.TicketNonce, credential.Expires)
		FailedToIssueTicket(request.Request, user.User, credential.TicketNonce, credential.Expires, error)
		IssueTicket(request.Request, user.User, credential.TicketNonce, credential.Expires)
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, credential.Expires)
		FailedToIssueApiToken(request.Request, user.User, credential.Expires, error)
		IssueApiToken(request.Request, user.User, credential.ApiRoles, credential.Expires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, credential.Expires)
		FailedToIssueContentToken(request.Request, user.User, credential.Expires, error)
		IssueContentToken(request.Request, user.User, credential.Expires)
	}
)
