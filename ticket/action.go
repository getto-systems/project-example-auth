package ticket

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		Register(request request.Request, user user.User, exp credential.Expiration) (credential.TicketNonce, credential.Expires, error)
		Validate(request request.Request, user user.User, ticket credential.Ticket) error
		Extend(request request.Request, user user.User, ticket credential.Ticket) (credential.Expires, error)
		Deactivate(request request.Request, user user.User, ticket credential.Ticket) error
	}
)
