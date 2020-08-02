package http_handler

import (
	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type PasswordLogin struct {
	Handler
}

func NewPasswordLogin(handler Handler) PasswordLogin {
	return PasswordLogin{Handler: handler}
}

func (handler PasswordLogin) handler() client.PasswordLoginHandler {
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
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}
