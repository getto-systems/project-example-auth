package ticket

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Register interface {
		Register(request request.Request, user user.User, exp Expiration) (credential.TicketNonce, time.Expires, error)
	}

	Validate interface {
		Validate(request request.Request, user user.User, ticket credential.Ticket) error
	}

	Extend interface {
		Extend(request request.Request, user user.User, ticket credential.Ticket) (time.Expires, error)
	}

	Deactivate interface {
		Deactivate(request request.Request, user user.User, ticket credential.Ticket) error
	}
)
