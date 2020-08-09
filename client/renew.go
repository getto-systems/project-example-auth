package client

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	Renew struct {
		Client
	}
	RenewHandler interface {
		RenewRequest() (request.Request, error)
		RenewResponse(error)
	}
)

func NewRenew(client Client) Renew {
	return Renew{
		Client: client,
	}
}

func (client Renew) Renew(handler RenewHandler) {
	credential, err := client.renew(handler)
	client.handleCredential(credential, err)
	handler.RenewResponse(err)
}
func (client Renew) renew(handler RenewHandler) (_ credential.Credential, err error) {
	ticket, err := client.getTicket()
	if err != nil {
		return
	}

	request, err := handler.RenewRequest()
	if err != nil {
		return
	}

	user, err := client.credential.ParseTicket(request, ticket)
	if err != nil {
		return
	}

	err = client.ticket.Validate(request, user, ticket)
	if err != nil {
		return
	}

	expires, err := client.ticket.Extend(request, user, ticket)
	if err != nil {
		return
	}

	return client.issueCredential(request, user, ticket.Nonce(), expires)
}
