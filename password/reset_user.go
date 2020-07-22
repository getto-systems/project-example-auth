package password

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrResetTokenUserNotMatched = errors.New("reset token user not matched")
	ErrResetTokenAlreadyExpired = errors.New("reset token already expired")
)

type (
	ResetUser struct {
		user    data.User
		expires data.Expires
	}
)

func NewResetUser(user data.User, expires data.Expires) ResetUser {
	return ResetUser{
		user:    user,
		expires: expires,
	}
}

func (reset ResetUser) User() data.User {
	return reset.user
}

func (reset ResetUser) Validate(request data.Request, loginUser data.User) error {
	if request.RequestedAt().Expired(reset.expires) {
		return ErrResetTokenAlreadyExpired
	}
	if reset.user.UserID() != loginUser.UserID() {
		return ErrResetTokenUserNotMatched
	}
	return nil
}
