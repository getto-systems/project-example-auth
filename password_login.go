package auth

import (
	"errors"

	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

var (
	ErrInvalidPasswordLogin = errors.New("invalid-password-login")
)

type (
	PasswordLogin struct {
		Usecase
	}
	PasswordLoginHandler interface {
		LoginRequest() (request.Request, user.Login, password.RawPassword, error)
		LoginResponse(error)
	}
)

func NewPasswordLogin(u Usecase) PasswordLogin {
	return PasswordLogin{Usecase: u}
}

func (u PasswordLogin) Login(handler PasswordLoginHandler) {
	handler.LoginResponse(u.login(handler))
}
func (u PasswordLogin) login(handler PasswordLoginHandler) (err error) {
	request, login, raw, err := handler.LoginRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	loginUser, err := u.user.GetUser(request, login)
	if err != nil {
		switch err {
		case user.ErrGetUserNotFoundUser:
			err = ErrInvalidPasswordLogin

		default:
			err = ErrServerError
		}
		return
	}

	extendSecond, err := u.password.Validate(request, loginUser, raw)
	if err != nil {
		switch err {
		case password.ErrCheckLengthEmpty,
			password.ErrCheckLengthTooLong,
			password.ErrValidateNotFoundPassword,
			password.ErrValidateMatchFailed:

			err = ErrInvalidPasswordLogin

		default:
			err = ErrServerError
		}
		return
	}

	newTicket, err := u.ticket.Register(request, loginUser, extendSecond)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	err = u.issueCredential(request, newTicket)
	if err != nil {
		return
	}

	return nil
}
