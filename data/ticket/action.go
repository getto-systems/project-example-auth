package ticket

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Action interface {
		Register(request request.Request, user user.User, exp Expiration) (credential.TicketNonce, time.Expires, error)
		Validate(request request.Request, user user.User, ticket credential.Ticket) error
		Extend(request request.Request, user user.User, ticket credential.Ticket) (time.Expires, error)
		Deactivate(request request.Request, user user.User, ticket credential.Ticket) error
	}
)
