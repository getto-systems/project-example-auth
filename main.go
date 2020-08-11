package _usecase

import (
	"errors"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Backend struct {
		ticket        ticket.Action
		credential    credential.Action
		user          user.Action
		password      password.Action
		passwordReset password_reset.Action
	}

	Usecase struct {
		Backend
		handler CredentialHandler
	}

	CredentialHandler interface {
		GetTicketNonceAndSignature() (credential.TicketNonce, credential.TicketSignature, error)
		SetCredential(credential.Credential)
		ClearCredential()
	}
)

func NewBackend(
	ticket ticket.Action,
	credential credential.Action,
	user user.Action,
	password password.Action,
	passwordReset password_reset.Action,
) Backend {
	return Backend{
		ticket:        ticket,
		credential:    credential,
		user:          user,
		password:      password,
		passwordReset: passwordReset,
	}
}

func NewUsecase(backend Backend, handler CredentialHandler) Usecase {
	return Usecase{
		Backend: backend,
		handler: handler,
	}
}

func (u Usecase) handleCredential(credential credential.Credential, err error) {
	if err != nil {
		u.handleCredentialError(err)
	} else {
		u.handler.SetCredential(credential)
	}
}
func (u Usecase) handleCredentialError(err error) {
	if errors.Is(err, ErrTicketValidate) {
		u.clearCredential()
	}
}
func (u Usecase) clearCredential() {
	u.handler.ClearCredential()
}
func (u Usecase) getTicketNonceAndSignature() (credential.TicketNonce, credential.TicketSignature, error) {
	return u.handler.GetTicketNonceAndSignature()
}
