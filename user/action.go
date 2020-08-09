package user

import (
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	Action interface {
		GetLogin(request request.Request, user User) (Login, error)
		GetUser(request request.Request, login Login) (User, error)
	}
)
