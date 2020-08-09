package credential

import (
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		ParseTicket(request request.Request, ticket Ticket) (user.User, error)
		IssueTicket(request request.Request, user user.User, nonce TicketNonce, expires Expires) (Ticket, error)
		IssueApiToken(request request.Request, user user.User, expires Expires) (ApiToken, error)
		IssueContentToken(request request.Request, user user.User, expires Expires) (ContentToken, error)
	}
)
