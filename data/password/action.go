package password

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Validate interface {
		Validate(request request.Request, user user.User, raw RawPassword) (ticket.Expiration, error)
	}

	Change interface {
		Change(request request.Request, user user.User, raw RawPassword) error
	}
)
