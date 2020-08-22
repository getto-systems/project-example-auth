package user

import (
	"github.com/getto-systems/project-example-auth/request"
)

type (
	Action interface {
		GetLogin(request.Request, User) (Login, error)
		GetUser(request.Request, Login) (User, error)
	}
)
