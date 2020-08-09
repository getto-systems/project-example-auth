package _usecase

import (
	"errors"

	"github.com/getto-systems/project-example-id/credential"
)

type (
	Usecase struct {
		Backend
		handler CredentialHandler
	}

	CredentialHandler interface {
		GetTicket() (credential.Ticket, error)
		SetCredential(credential.Credential)
		ClearCredential()
	}
)

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
func (u Usecase) getTicket() (credential.Ticket, error) {
	return u.handler.GetTicket()
}
