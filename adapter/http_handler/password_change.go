package http_handler

import (
	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
)

type PasswordChange struct {
	Handler
}

func NewPasswordChange(handler Handler) PasswordChange {
	return PasswordChange{Handler: handler}
}

func (handler PasswordChange) handler() client.PasswordChangeHandler {
	return handler
}

func (handler PasswordChange) GetLoginRequest() (request.Request, error) {
	return handler.newRequest(), nil
}
func (handler PasswordChange) GetLoginResponse(login user.Login, err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	type body struct {
		LoginID string `json:"login_id"`
	}

	handler.ok(body{
		LoginID: string(login.ID()),
	})
}

func (handler PasswordChange) ChangeRequest() (_ request.Request, _ password.ChangeParam, err error) {
	type body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var input body
	err = handler.parseBody(&input)
	if err != nil {
		return
	}

	param := password.ChangeParam{
		OldPassword: password.RawPassword(input.OldPassword),
		NewPassword: password.RawPassword(input.NewPassword),
	}

	return handler.newRequest(), param, nil
}
func (handler PasswordChange) ChangeResponse(err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}
