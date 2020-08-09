package infra

import (
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger interface {
		ValidateLogger
		ChangeLogger
	}

	ValidateLogger interface {
		TryToValidate(request.Request, user.User)
		FailedToValidate(request.Request, user.User, error)
		FailedToValidateBecausePasswordCheckFailed(request.Request, user.User, error)
		FailedToValidateBecausePasswordNotFound(request.Request, user.User, error)
		FailedToValidateBecausePasswordMatchFailed(request.Request, user.User, error)
		AuthByPassword(request.Request, user.User)
	}

	ChangeLogger interface {
		TryToChange(request.Request, user.User)
		FailedToChange(request.Request, user.User, error)
		FailedToChangeBecausePasswordCheckFailed(request.Request, user.User, error)
		Change(request.Request, user.User)
	}
)
