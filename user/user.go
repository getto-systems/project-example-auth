package user

import (
	"errors"
)

var (
	ErrGetLoginNotFoundLogin = errors.New("User.GetLogin/NotFound.Login")

	ErrGetUserNotFoundUser = errors.New("User.GetUser/NotFound.User")
)

type (
	UserID  string
	LoginID string

	User struct {
		id UserID
	}
	UserLog struct {
		ID string `json:"id"`
	}

	Login struct {
		id LoginID
	}
	LoginLog struct {
		ID string `json:"id"`
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

func (user User) Log() UserLog {
	return UserLog{
		ID: string(user.ID()),
	}
}
func (login Login) Log() LoginLog {
	return LoginLog{
		ID: string(login.ID()),
	}
}
