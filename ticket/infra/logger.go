package infra

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
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
		Register(request.Request, user.User, expiration.Expires, expiration.ExtendLimit, credential.TicketNonce)
	}

	ValidateLogger interface {
		TryToValidate(request.Request, user.User, credential.TicketNonce)
		FailedToValidate(request.Request, user.User, credential.TicketNonce, error)
		FailedToValidateBecauseExpired(request.Request, user.User, credential.TicketNonce, error)
		FailedToValidateBecauseTicketNotFound(request.Request, user.User, credential.TicketNonce, error)
		FailedToValidateBecauseUserMatchFailed(request.Request, user.User, credential.TicketNonce, error)
		AuthByTicket(request.Request, user.User, credential.TicketNonce)
	}

	DeactivateLogger interface {
		TryToDeactivate(request.Request, user.User, credential.TicketNonce)
		FailedToDeactivate(request.Request, user.User, credential.TicketNonce, error)
		Deactivate(request.Request, user.User, credential.TicketNonce)
	}

	ExtendLogger interface {
		TryToExtend(request.Request, user.User, credential.TicketNonce)
		FailedToExtend(request.Request, user.User, credential.TicketNonce, error)
		Extend(request.Request, user.User, credential.TicketNonce, expiration.Expires, expiration.ExtendLimit)
	}
)
