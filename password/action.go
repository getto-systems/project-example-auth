package password

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	Action interface {
		Validate(request.Request, user.User, RawPassword) (credential.TicketExtendSecond, error)
		Change(request.Request, user.User, RawPassword) error
	}
)
