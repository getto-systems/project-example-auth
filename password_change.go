package auth

import (
	"errors"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	ErrInvalidOldPassword = errors.New("invalid-old-password")
	ErrInvalidNewPassword = errors.New("invalid-new-password")
)

type (
	PasswordChange struct {
		Usecase
	}
	PasswordChangeHandler interface {
		GetLoginRequest() (request.Request, error)
		GetLoginResponse(user.Login, error)

		ChangeRequest() (request.Request, password.ChangeParam, error)
		ChangeResponse(error)
	}
)

func NewPasswordChange(u Usecase) PasswordChange {
	return PasswordChange{Usecase: u}
}

func (u PasswordChange) GetLogin(handler PasswordChangeHandler) {
	handler.GetLoginResponse(u.getLogin(handler))
}
func (u PasswordChange) getLogin(handler PasswordChangeHandler) (_ user.Login, err error) {
	request, err := handler.GetLoginRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	user, _, err := u.validateTicket(request)
	if err != nil {
		return
	}

	login, err := u.user.GetLogin(request, user)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	return login, nil
}

func (u PasswordChange) Change(handler PasswordChangeHandler) {
	handler.ChangeResponse(u.change(handler))
}
func (u PasswordChange) change(handler PasswordChangeHandler) (err error) {
	request, param, err := handler.ChangeRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	user, _, err := u.validateTicket(request)
	if err != nil {
		return
	}

	_, err = u.password.Validate(request, user, param.OldPassword)
	if err != nil {
		switch err {
		case password.ErrCheckLengthEmpty,
			password.ErrCheckLengthTooLong,
			password.ErrValidateMatchFailed,
			password.ErrValidateNotFoundPassword:

			err = ErrInvalidOldPassword

		default:
			err = ErrServerError
		}
		return
	}

	err = u.password.Change(request, user, param.NewPassword)
	if err != nil {
		switch err {
		case password.ErrCheckLengthEmpty,
			password.ErrCheckLengthTooLong:

			err = ErrInvalidNewPassword

		default:
			err = ErrServerError
		}
		return
	}

	return nil
}
