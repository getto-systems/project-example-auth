package infra

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger interface {
		RegisterLogger
		ValidateLogger
		DeactivateLogger
		ExtendLogger
	}

	RegisterLogger interface {
		TryToRegister(request.Request, user.User, expiration.Expires, expiration.ExtendLimit)
		FailedToRegister(request.Request, user.User, expiration.Expires, expiration.ExtendLimit, error)
		Register(request.Request, user.User, expiration.Expires, expiration.ExtendLimit)
	}

	ValidateLogger interface {
		TryToValidate(request.Request, user.User)
		FailedToValidate(request.Request, user.User, error)
		FailedToValidateBecauseExpired(request.Request, user.User, error)
		FailedToValidateBecauseTicketNotFound(request.Request, user.User, error)
		FailedToValidateBecauseUserMatchFailed(request.Request, user.User, error)
		AuthByTicket(request.Request, user.User)
	}

	DeactivateLogger interface {
		TryToDeactivate(request.Request, user.User)
		FailedToDeactivate(request.Request, user.User, error)
		Deactivate(request.Request, user.User)
	}

	ExtendLogger interface {
		TryToExtend(request.Request, user.User)
		FailedToExtend(request.Request, user.User, error)
		Extend(request.Request, user.User, expiration.Expires, expiration.ExtendLimit)
	}
)
