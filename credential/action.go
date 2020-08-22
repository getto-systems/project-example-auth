package credential

import (
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	Action interface {
		ParseTicketSignature(request.Request, TicketNonce, TicketSignature) (user.User, error)
		IssueTicketToken(request.Request, Ticket) (TicketToken, error)
		IssueApiToken(request.Request, Ticket) (ApiToken, error)
		IssueContentToken(request.Request, Ticket) (ContentToken, error)
	}
)
