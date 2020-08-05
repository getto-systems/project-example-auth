package client

import (
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	Logout struct {
		Client
	}
	LogoutHandler interface {
		LogoutRequest() (request.Request, error)
		LogoutResponse(error)
	}
)

func NewLogout(client Client) Logout {
	return Logout{Client: client}
}

func (client Logout) Logout(handler LogoutHandler) {
	err := client.logout(handler)
	client.clearCredential()
	handler.LogoutResponse(err)
}
func (client Logout) logout(handler LogoutHandler) (err error) {
	ticket, err := client.getTicket()
	if err != nil {
		return
	}

	request, err := handler.LogoutRequest()
	if err != nil {
		return
	}

	user, err := client.ticket.validate.Validate(request, ticket)
	if err != nil {
		return
	}

	err = client.ticket.deactivate.Deactivate(request, user, ticket)
	if err != nil {
		return
	}

	return nil
}
