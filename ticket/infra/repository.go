package infra

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(credential.TicketNonce) (user.User, credential.TicketExpires, bool, error)

		FindExtendLimit(credential.TicketNonce) (credential.TicketExtendLimit, bool, error)
		UpdateExpires(credential.TicketNonce, credential.TicketExpires) error

		FindUser(credential.TicketNonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(credential.TicketNonce) error

		RegisterTicket(TicketNonceGenerator, user.User, credential.TicketExpires, credential.TicketExtendLimit) (credential.TicketNonce, error)
	}
)
