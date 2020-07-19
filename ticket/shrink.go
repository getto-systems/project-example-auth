package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ShrinkEventPublisher interface {
	ShrinkTicket(data.Request, Nonce, data.User)
	ShrinkTicketFailed(data.Request, Nonce, data.User, error)
}

type ShrinkDB interface {
	TicketExists(Nonce, data.User) bool
	ShrinkTicket(Nonce) error
}
