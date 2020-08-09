package ticket

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(api_token.TicketNonce) (user.User, time.Expires, bool, error)

		FindExpireSecondAndExtendLimit(api_token.TicketNonce) (time.Second, time.ExtendLimit, bool, error)

		FindUser(api_token.TicketNonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(api_token.TicketNonce) error

		RegisterTicket(api_token.TicketNonceGenerator, user.User, time.Expires, time.Second, time.ExtendLimit) (api_token.TicketNonce, error)
	}
)
