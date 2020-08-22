package http_handler

import (
	"github.com/getto-systems/project-example-auth"

	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type PasswordLogin struct {
	Handler
}

func NewPasswordLogin(handler Handler) PasswordLogin {
	return PasswordLogin{Handler: handler}
}

func (handler PasswordLogin) handler() auth.PasswordLoginHandler {
	return handler
}

func (handler PasswordLogin) LoginRequest() (_ request.Request, _ user.Login, _ password.RawPassword, err error) {
	type body struct {
		LoginID  string `json:"login_id"`
		Password string `json:"password"`
	}

	var input body
	err = handler.parseBody(&input)
	if err != nil {
		return
	}

	login := user.NewLogin(user.LoginID(input.LoginID))
	raw := password.RawPassword(input.Password)

	return handler.newRequest(), login, raw, nil
}
func (handler PasswordLogin) LoginResponse(err error) {
	if err != nil {
		switch err {
		case auth.ErrBadRequest:
			handler.badRequest()

		case auth.ErrInvalidPasswordLogin:
			handler.unauthorized("invalid-password-login")

		default:
			handler.internalServerError()
		}
		return
	}

	handler.ok("OK")
}
