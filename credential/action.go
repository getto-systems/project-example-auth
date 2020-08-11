package credential

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		ParseTicket(request.Request, TicketToken) (user.User, error)
		IssueTicket(request.Request, user.User, TicketNonce, expiration.Expires) (TicketToken, error)
		IssueApiToken(request.Request, user.User, expiration.Expires) (ApiToken, error)
		IssueContentToken(request.Request, user.User, expiration.Expires) (ContentToken, error)
	}
)
