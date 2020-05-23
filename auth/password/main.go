package password

import (
	"github.com/getto-systems/project-example-id/auth"
)

type UserPassword struct {
	db     UserPasswordRepository
	userID auth.UserID
}

type Password string

type UserPasswordRepository interface {
	MatchUserPassword(auth.UserID, Password) bool
}

func NewUserPassword(db UserPasswordRepository, userID auth.UserID) UserPassword {
	return UserPassword{db, userID}
}

func (userPassword UserPassword) Match(password Password) bool {
	return userPassword.db.MatchUserPassword(userPassword.userID, password)
}
