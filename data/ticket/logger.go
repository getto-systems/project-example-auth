package ticket

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Logger interface {
		ValidateLogger
		ShrinkLogger
		IssueLogger
		ExtendLogger
	}

	ValidateLogger interface {
		TryToValidateTicket(request.Request, Nonce)
		FailedToValidateTicket(request.Request, Nonce, error)
		FailedToValidateTicketBecauseExpired(request.Request, Nonce, error)
		AuthByTicket(request.Request, user.User, Nonce)
	}

	ShrinkLogger interface {
		TryToShrinkTicket(request.Request, user.User, Nonce)
		FailedToShrinkTicket(request.Request, user.User, Nonce, error)
		ShrinkTicket(request.Request, user.User, Nonce)
	}

	IssueLogger interface {
		TryToIssue(request.Request, user.User, time.Expires, time.ExtendLimit)
		FailedToIssue(request.Request, user.User, time.Expires, time.ExtendLimit, error)
		Issue(request.Request, user.User, time.Expires, time.ExtendLimit, Nonce)
	}

	ExtendLogger interface {
		TryToExtendTicket(request.Request, user.User, Nonce, time.Expires)
		FailedToExtendTicket(request.Request, user.User, Nonce, time.Expires, error)
		ExtendTicket(request.Request, user.User, Nonce, time.Expires)
	}
)
