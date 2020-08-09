package credential

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
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
		TryToIssueTicket(request.Request, user.User, credential.TicketNonce, time.Expires)
		FailedToIssueTicket(request.Request, user.User, credential.TicketNonce, time.Expires, error)
		IssueTicket(request.Request, user.User, credential.TicketNonce, time.Expires)
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, time.Expires)
		FailedToIssueApiToken(request.Request, user.User, time.Expires, error)
		IssueApiToken(request.Request, user.User, credential.ApiRoles, time.Expires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, time.Expires)
		FailedToIssueContentToken(request.Request, user.User, time.Expires, error)
		IssueContentToken(request.Request, user.User, time.Expires)
	}
)
