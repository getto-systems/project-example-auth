package infra

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

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
		TryToIssueTicket(request.Request, user.User, credential.TicketNonce, expiration.Expires)
		FailedToIssueTicket(request.Request, user.User, credential.TicketNonce, expiration.Expires, error)
		IssueTicket(request.Request, user.User, credential.TicketNonce, expiration.Expires)
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, expiration.Expires)
		FailedToIssueApiToken(request.Request, user.User, expiration.Expires, error)
		IssueApiToken(request.Request, user.User, credential.ApiRoles, expiration.Expires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, expiration.Expires)
		FailedToIssueContentToken(request.Request, user.User, expiration.Expires, error)
		IssueContentToken(request.Request, user.User, expiration.Expires)
	}
)
