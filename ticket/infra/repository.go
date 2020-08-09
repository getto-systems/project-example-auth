package infra

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(credential.TicketNonce) (user.User, expiration.Expires, bool, error)

		FindExtendLimit(credential.TicketNonce) (expiration.ExtendLimit, bool, error)
		UpdateExpires(credential.TicketNonce, expiration.Expires) error

		FindUser(credential.TicketNonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(credential.TicketNonce) error

		RegisterTicket(TicketNonceGenerator, user.User, expiration.Expires, expiration.ExtendLimit) (credential.TicketNonce, error)
	}
)
