package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

var (
	ticketExpireDuration = data.Minute(30)
)

type TicketSign interface {
	Parse(data.SignedTicket) (data.Ticket, error)
	Sign(data.Ticket) (data.SignedTicket, error)
}

func NewTicket(request data.Request, profile data.Profile) data.Ticket {
	requestedAt := request.RequestedAt
	expires := requestedAt.Expires(ticketExpireDuration)
	authenticatedAt := data.AuthenticatedAt(requestedAt)

	return data.Ticket{
		Profile:         profile,
		AuthenticatedAt: authenticatedAt,
		Expires:         expires,
	}
}
