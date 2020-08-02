package client

import (
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	PasswordChange struct {
		Client
	}
	PasswordChangeHandler interface {
		GetLoginRequest() (request.Request, error)
		GetLoginResponse(user.Login, error)

		ChangeRequest() (request.Request, password.ChangeParam, error)
		ChangeResponse(error)
	}
)

func NewPasswordChange(client Client) PasswordChange {
	return PasswordChange{Client: client}
}

func (client PasswordChange) GetLogin(handler PasswordChangeHandler) {
	login, err := client.getLogin(handler)
	client.handleCredentialError(err)
	handler.GetLoginResponse(login, err)
}
func (client PasswordChange) getLogin(handler PasswordChangeHandler) (_ user.Login, err error) {
	ticket, err := client.getTicket()
	if err != nil {
		return
	}

	request, err := handler.GetLoginRequest()
	if err != nil {
		return
	}

	user, err := client.ticket.validate.Validate(request, ticket)
	if err != nil {
		return
	}

	login, err := client.user.getLogin.Get(request, user)
	if err != nil {
		return
	}

	return login, nil
}

func (client PasswordChange) Change(handler PasswordChangeHandler) {
	err := client.change(handler)
	client.handleCredentialError(err)
	handler.ChangeResponse(err)
}
func (client PasswordChange) change(handler PasswordChangeHandler) (err error) {
	ticket, err := client.getTicket()
	if err != nil {
		return
	}

	request, param, err := handler.ChangeRequest()
	if err != nil {
		return
	}

	user, err := client.ticket.validate.Validate(request, ticket)
	if err != nil {
		return
	}

	err = client.password.validate.Validate(request, user, param.OldPassword)
	if err != nil {
		return
	}

	err = client.password.change.Change(request, user, param.NewPassword)
	if err != nil {
		return
	}

	return nil
}
