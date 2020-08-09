package client

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/credential"
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
	if errors.Is(err, data.ErrTicketValidate) {
		client.clearCredential()
	}
}
func (client Client) clearCredential() {
	client.handler.ClearCredential()
}
func (client Client) getTicket() (credential.Ticket, error) {
	return client.handler.GetTicket()
}
