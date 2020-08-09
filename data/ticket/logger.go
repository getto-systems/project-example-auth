package ticket

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Logger interface {
		ValidateLogger
		DeactivateLogger
		IssueLogger
		ExtendLogger
	}

	ValidateLogger interface {
		TryToValidate(request.Request, credential.TicketNonce)
		FailedToValidate(request.Request, credential.TicketNonce, error)
		FailedToValidateBecauseExpired(request.Request, credential.TicketNonce, error)
		FailedToValidateBecauseTicketNotFound(request.Request, credential.TicketNonce, error)
		FailedToValidateBecauseMatchFailed(request.Request, credential.TicketNonce, error)
		AuthByTicket(request.Request, user.User, credential.TicketNonce)
	}

	DeactivateLogger interface {
		TryToDeactivate(request.Request, user.User, credential.TicketNonce)
		FailedToDeactivate(request.Request, user.User, credential.TicketNonce, error)
		Deactivate(request.Request, user.User, credential.TicketNonce)
	}

	IssueLogger interface {
		TryToIssue(request.Request, user.User, time.Expires, time.ExtendLimit)
		FailedToIssue(request.Request, user.User, time.Expires, time.ExtendLimit, error)
		Issue(request.Request, user.User, time.Expires, time.ExtendLimit, credential.TicketNonce)
	}

	ExtendLogger interface {
		TryToExtend(request.Request, user.User, credential.TicketNonce)
		FailedToExtend(request.Request, user.User, credential.TicketNonce, error)
		Extend(request.Request, user.User, credential.TicketNonce, time.Expires, time.ExtendLimit)
	}
)
