package user

import (
	"github.com/getto-systems/project-example-id/_misc/errors"
)

var (
	ErrGetLoginNotFoundLogin = errors.NewError("User.GetLogin", "NotFound.Login")

	ErrGetUserNotFoundUser = errors.NewError("User.GetUser", "NotFound.User")
)

type (
	User struct {
		id UserID
	}
	UserID string

	LoginID string
	Login   struct {
		id LoginID
	}
)

func NewUser(userID UserID) User {
	return User{
		id: userID,
	}
}

func (user User) ID() UserID {
	return user.id
}

func NewLogin(loginID LoginID) Login {
	return Login{
		id: loginID,
	}
}

func (login Login) ID() LoginID {
	return login.id
}
