package ticket

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketRepository interface {
		FindUserAndExtendLimit(Nonce) (user.User, time.ExtendLimit, bool, error)

		FindUser(Nonce) (user.User, bool, error)
		ShrinkExtendLimit(Nonce) error

		RegisterTicket(NonceGenerator, user.User, time.Expires, time.ExtendLimit) (Nonce, error)
	}
)
