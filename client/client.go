package client

import (
	goerrors "errors"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/ticket"
)

type (
	Client struct {
		Backend
		handler CredentialHandler
	}

	CredentialHandler interface {
		GetTicket() (ticket.Ticket, error)
		SetCredential(data.Credential)
		ClearCredential()
	}
)

func NewClient(backend Backend, handler CredentialHandler) Client {
	return Client{
		Backend: backend,
		handler: handler,
	}
}

func (client Client) handleCredential(credential data.Credential, err error) {
	if err != nil {
		client.handleCredentialError(err)
	} else {
		client.handler.SetCredential(credential)
	}
}
func (client Client) handleCredentialError(err error) {
	var e *errors.Error
	if goerrors.As(err, &e) {
		if e.TicketValidateError() {
			client.handler.ClearCredential()
		}
	}
}
func (client Client) clearCredential() {
	client.handler.ClearCredential()
}
func (client Client) getTicket() (ticket.Ticket, error) {
	return client.handler.GetTicket()
}
