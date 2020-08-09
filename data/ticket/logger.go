package ticket

import (
	"github.com/getto-systems/project-example-id/data/api_token"
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
		TryToValidate(request.Request, api_token.TicketNonce)
		FailedToValidate(request.Request, api_token.TicketNonce, error)
		FailedToValidateBecauseExpired(request.Request, api_token.TicketNonce, error)
		FailedToValidateBecauseTicketNotFound(request.Request, api_token.TicketNonce, error)
		FailedToValidateBecauseMatchFailed(request.Request, api_token.TicketNonce, error)
		AuthByTicket(request.Request, user.User, api_token.TicketNonce)
	}

	DeactivateLogger interface {
		TryToDeactivate(request.Request, user.User, api_token.TicketNonce)
		FailedToDeactivate(request.Request, user.User, api_token.TicketNonce, error)
		Deactivate(request.Request, user.User, api_token.TicketNonce)
	}

	IssueLogger interface {
		TryToIssue(request.Request, user.User, time.Expires, time.ExtendLimit)
		FailedToIssue(request.Request, user.User, time.Expires, time.ExtendLimit, error)
		Issue(request.Request, user.User, time.Expires, time.ExtendLimit, api_token.TicketNonce)
	}

	ExtendLogger interface {
		TryToExtend(request.Request, user.User, api_token.TicketNonce)
		FailedToExtend(request.Request, user.User, api_token.TicketNonce, error)
		Extend(request.Request, user.User, api_token.TicketNonce, time.Expires, time.ExtendLimit)
	}
)
