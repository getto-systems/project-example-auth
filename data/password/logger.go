package password

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Logger interface {
		ValidateLogger
		ChangeLogger
	}

	ValidateLogger interface {
		TryToValidate(request.Request, user.User)
		FailedToValidate(request.Request, user.User, error)
		AuthByPassword(request.Request, user.User)
	}

	ChangeLogger interface {
		TryToChange(request.Request, user.User)
		FailedToChange(request.Request, user.User, error)
		Change(request.Request, user.User)
	}
)