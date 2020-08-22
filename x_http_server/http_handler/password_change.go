package http_handler

import (
	"github.com/getto-systems/project-example-auth"

	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type PasswordChange struct {
	Handler
}

func NewPasswordChange(handler Handler) PasswordChange {
	return PasswordChange{Handler: handler}
}

func (handler PasswordChange) handler() auth.PasswordChangeHandler {
	return handler
}

func (handler PasswordChange) GetLoginRequest() (request.Request, error) {
	return handler.newRequest(), nil
}
func (handler PasswordChange) GetLoginResponse(login user.Login, err error) {
	if err != nil {
		switch err {
		case auth.ErrBadRequest:
			handler.badRequest()

		case auth.ErrInvalidTicket:
			handler.invalidTicket()

		default:
			handler.internalServerError()
		}
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
		switch err {
		case auth.ErrBadRequest:
			handler.badRequest()

		case auth.ErrInvalidTicket:
			handler.invalidTicket()

		case auth.ErrInvalidOldPassword:
			handler.unauthorized("invalid-old-password")

		default:
			handler.internalServerError()
		}
		return
	}

	handler.ok("OK")
}
