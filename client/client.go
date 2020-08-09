package client

import (
	goerrors "errors"

	"github.com/getto-systems/project-example-id/misc/errors"

	"github.com/getto-systems/project-example-id/credential"
)

var (
	ErrTicketValidate = errors.NewError("Ticket.Validate", "")
	ErrPasswordCheck  = errors.NewError("Password.Check", "")
)

type (
	Client struct {
		Backend
		handler CredentialHandler
	}

	CredentialHandler interface {
		GetTicket() (credential.Ticket, error)
		SetCredential(credential.Credential)
		ClearCredential()
	}
)

func NewClient(backend Backend, handler CredentialHandler) Client {
	return Client{
		Backend: backend,
		handler: handler,
	}
}

func (client Client) handleCredential(credential credential.Credential, err error) {
	if err != nil {
		client.handleCredentialError(err)
	} else {
		client.handler.SetCredential(credential)
	}
}
func (client Client) handleCredentialError(err error) {
	if goerrors.Is(err, ErrTicketValidate) {
		client.clearCredential()
	}
}
func (client Client) clearCredential() {
	client.handler.ClearCredential()
}
func (client Client) getTicket() (credential.Ticket, error) {
	return client.handler.GetTicket()
}
