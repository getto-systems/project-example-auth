package _usecase

import (
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
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
	login, err := u.getLogin(handler)
	u.handleCredentialError(err)
	handler.GetLoginResponse(login, err)
}
func (u PasswordChange) getLogin(handler PasswordChangeHandler) (_ user.Login, err error) {
	nonce, signature, err := u.getTicketNonceAndSignature()
	if err != nil {
		return
	}

	request, err := handler.GetLoginRequest()
	if err != nil {
		return
	}

	user, err := u.credential.ParseTicketSignature(request, nonce, signature)
	if err != nil {
		return
	}

	err = u.ticket.Validate(request, user, nonce)
	if err != nil {
		return
	}

	login, err := u.user.GetLogin(request, user)
	if err != nil {
		return
	}

	return login, nil
}

func (u PasswordChange) Change(handler PasswordChangeHandler) {
	err := u.change(handler)
	u.handleCredentialError(err)
	handler.ChangeResponse(err)
}
func (u PasswordChange) change(handler PasswordChangeHandler) (err error) {
	nonce, signature, err := u.getTicketNonceAndSignature()
	if err != nil {
		return
	}

	request, param, err := handler.ChangeRequest()
	if err != nil {
		return
	}

	user, err := u.credential.ParseTicketSignature(request, nonce, signature)
	if err != nil {
		return
	}

	err = u.ticket.Validate(request, user, nonce)
	if err != nil {
		return
	}

	_, err = u.password.Validate(request, user, param.OldPassword)
	if err != nil {
		return
	}

	err = u.password.Change(request, user, param.NewPassword)
	if err != nil {
		return
	}

	return nil
}
