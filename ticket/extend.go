package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ExtendEventPublisher interface {
	ExtendTicket(data.Request, Nonce, data.User, data.Expires)
	ExtendTicketFailed(data.Request, Nonce, data.User, data.Expires, error)
}

type ExtendDB interface {
	FindTicketExtendLimit(Nonce, data.User) (data.ExtendLimit, error)
}
