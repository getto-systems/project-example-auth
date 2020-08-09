package user

import (
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	GetLogin interface {
		Get(request request.Request, user User) (Login, error)
	}

	GetUser interface {
		Get(request request.Request, login Login) (User, error)
	}
)
