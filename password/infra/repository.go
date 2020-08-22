package infra

import (
	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	PasswordRepository interface {
		FindPassword(user.User) (password.HashedPassword, bool, error)

		ChangePassword(user.User, password.HashedPassword) error
	}
)
