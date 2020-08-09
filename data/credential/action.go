package credential

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	ParseTicket interface {
		Parse(request request.Request, ticket Ticket) (user.User, error)
	}

	IssueTicket interface {
		Issue(request request.Request, user user.User, nonce TicketNonce, expires time.Expires) (Ticket, error)
	}

	IssueApiToken interface {
		Issue(request request.Request, user user.User, expires time.Expires) (ApiToken, error)
	}

	IssueContentToken interface {
		Issue(request request.Request, user user.User, expires time.Expires) (ContentToken, error)
	}
)
