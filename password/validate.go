package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type ValidateLogger interface {
	TryToValidate(data.Request, Login)
	FailedToValidate(data.Request, Login, error)
	AuthedByPassword(data.Request, Login, data.User)
}
