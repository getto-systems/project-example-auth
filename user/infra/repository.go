package infra

import (
	"github.com/getto-systems/project-example-id/user"
)

type (
	UserRepository interface {
		FindUser(user.Login) (user.User, bool, error)
		FindLogin(user.User) (user.Login, bool, error)
		RegisterUser(user.User, user.Login) error
	}
)
