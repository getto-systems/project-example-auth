package ticket

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		Register(request.Request, user.User, expiration.ExtendSecond) (credential.TicketNonce, expiration.Expires, error)
		Validate(request.Request, user.User, credential.TicketToken) error
		Extend(request.Request, user.User, credential.TicketToken) (expiration.Expires, error)
		Deactivate(request.Request, user.User, credential.TicketToken) error
	}
)
