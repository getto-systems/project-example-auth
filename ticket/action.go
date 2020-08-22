package ticket

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	Action interface {
		Register(request.Request, user.User, credential.TicketExtendSecond) (credential.Ticket, error)
		Validate(request.Request, user.User, credential.TicketNonce) error
		Extend(request.Request, user.User, credential.TicketNonce) (credential.Ticket, error)
		Deactivate(request.Request, user.User, credential.TicketNonce) error
	}
)
