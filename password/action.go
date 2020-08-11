package password

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		Validate(request.Request, user.User, RawPassword) (credential.TicketExtendSecond, error)
		Change(request.Request, user.User, RawPassword) error
	}
)
