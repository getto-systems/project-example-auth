package password

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Action interface {
		Validate(request request.Request, user user.User, raw RawPassword) (ticket.Expiration, error)
		Change(request request.Request, user user.User, raw RawPassword) error
	}
)
