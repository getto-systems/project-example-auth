package _usecase

import (
	"github.com/getto-systems/project-example-id/request"
)

type (
	Logout struct {
		Usecase
	}
	LogoutHandler interface {
		LogoutRequest() (request.Request, error)
		LogoutResponse(error)
	}
)

func NewLogout(u Usecase) Logout {
	return Logout{Usecase: u}
}

func (u Logout) Logout(handler LogoutHandler) {
	err := u.logout(handler)
	u.clearCredential()
	handler.LogoutResponse(err)
}
func (u Logout) logout(handler LogoutHandler) (err error) {
	ticket, err := u.getTicket()
	if err != nil {
		return
	}

	request, err := handler.LogoutRequest()
	if err != nil {
		return
	}

	user, err := u.credential.ParseTicket(request, ticket)
	if err != nil {
		return
	}

	err = u.ticket.Validate(request, user, ticket)
	if err != nil {
		return
	}

	err = u.ticket.Deactivate(request, user, ticket)
	if err != nil {
		return
	}

	return nil
}
