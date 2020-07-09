package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type TicketSign interface {
	Parse(data.SignedTicket) (data.Ticket, error)
	Sign(data.Ticket) (data.SignedTicket, error)
}
