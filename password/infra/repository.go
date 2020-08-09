package infra

import (
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
)

type (
	PasswordRepository interface {
		FindPassword(user.User) (password.HashedPassword, bool, error)

		ChangePassword(user.User, password.HashedPassword) error
	}
)
