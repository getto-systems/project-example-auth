package infra

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger interface {
		GetLoginLogger
		GetUserLogger
	}

	GetLoginLogger interface {
		TryToGetLogin(request.Request, user.User)
		FailedToGetLogin(request.Request, user.User, error)
		GetLogin(request.Request, user.User, user.Login)
	}

	GetUserLogger interface {
		TryToGetUser(request.Request, user.Login)
		FailedToGetUser(request.Request, user.Login, error)
		FailedToGetUserBecauseUserNotFound(request.Request, user.Login, error)
		GetUser(request.Request, user.Login, user.User)
	}
)
