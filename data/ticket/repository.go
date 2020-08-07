package ticket

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(Nonce) (user.User, time.Expires, bool, error)

		FindExpireSecondAndExtendLimit(Nonce) (time.Second, time.ExtendLimit, bool, error)

		FindUser(Nonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(Nonce) error

		RegisterTicket(NonceGenerator, user.User, time.Expires, time.Second, time.ExtendLimit) (Nonce, error)
	}
)
