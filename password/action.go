package password

import (
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		Validate(request request.Request, user user.User, raw RawPassword) (ticket.Expiration, error)
		Change(request request.Request, user user.User, raw RawPassword) error
	}
)
