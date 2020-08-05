package user

import (
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	Logger interface {
		GetLoginLogger
		GetUserLogger
	}

	GetLoginLogger interface {
		TryToGetLogin(request.Request, User)
		FailedToGetLogin(request.Request, User, error)
		GetLogin(request.Request, User, Login)
	}

	GetUserLogger interface {
		TryToGetUser(request.Request, Login)
		FailedToGetUser(request.Request, Login, error)
		FailedToGetUserBecauseUserNotFound(request.Request, Login, error)
		GetUser(request.Request, Login, User)
	}
)
