package password

import (
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	PasswordRepository interface {
		FindPassword(user.User) (password.HashedPassword, bool, error)

		ChangePassword(user.User, password.HashedPassword) error
	}
)
