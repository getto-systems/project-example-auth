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
		TryToValidate(request.Request, Nonce)
		FailedToValidate(request.Request, Nonce, error)
		FailedToValidateBecauseExpired(request.Request, Nonce, error)
		AuthByTicket(request.Request, user.User, Nonce)
	}

	ShrinkLogger interface {
		TryToShrink(request.Request, user.User, Nonce)
		FailedToShrink(request.Request, user.User, Nonce, error)
		Shrink(request.Request, user.User, Nonce)
	}

	IssueLogger interface {
		TryToIssue(request.Request, user.User, time.Expires, time.ExtendLimit)
		FailedToIssue(request.Request, user.User, time.Expires, time.ExtendLimit, error)
		Issue(request.Request, user.User, time.Expires, time.ExtendLimit, Nonce)
	}

	ExtendLogger interface {
		TryToExtend(request.Request, user.User, Nonce, time.Expires)
		FailedToExtend(request.Request, user.User, Nonce, time.Expires, error)
		Extend(request.Request, user.User, Nonce, time.Expires)
	}
)