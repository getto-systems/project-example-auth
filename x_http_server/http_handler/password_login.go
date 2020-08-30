package http_handler

import (
	"log"

	"github.com/getto-systems/project-example-auth/y_static/protocol_buffers/password_login_pb"

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
	var passwordLogin password_login_pb.PasswordLoginMessage
	err = handler.parseBodyProto(&passwordLogin)
	if err != nil {
		log.Printf("body parse failed: %s", err)
		return
	}

	login := user.NewLogin(user.LoginID(passwordLogin.LoginId))
	raw := password.RawPassword(passwordLogin.Password)

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
