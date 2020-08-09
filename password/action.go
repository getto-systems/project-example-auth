package password

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		Validate(request.Request, user.User, RawPassword) (expiration.ExtendSecond, error)
		Change(request.Request, user.User, RawPassword) error
	}
)
