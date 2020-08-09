package infra

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(credential.TicketNonce) (user.User, credential.Expires, bool, error)

		FindExpireSecondAndExtendLimit(credential.TicketNonce) (time.Second, credential.ExtendLimit, bool, error)
		UpdateExpires(credential.TicketNonce, credential.Expires) error

		FindUser(credential.TicketNonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(credential.TicketNonce) error

		RegisterTicket(TicketNonceGenerator, user.User, credential.Expires, time.Second, credential.ExtendLimit) (credential.TicketNonce, error)
	}
)
