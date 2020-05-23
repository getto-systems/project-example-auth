package password

import (
	"github.com/getto-systems/project-example-id/user"
)

type UserPassword string

type UserPasswordRepository interface {
	MatchUserPassword(user.UserID, UserPassword) bool
}
