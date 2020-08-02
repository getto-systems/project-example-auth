package password

import (
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	PasswordRepository interface {
		FindPassword(user.User) (HashedPassword, bool, error)

		ChangePassword(user.User, HashedPassword) error
	}
)
