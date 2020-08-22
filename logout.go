package auth

import (
	"github.com/getto-systems/project-example-auth/request"
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
	handler.LogoutResponse(u.logout(handler))
}
func (u Logout) logout(handler LogoutHandler) (err error) {
	request, err := handler.LogoutRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	user, nonce, err := u.validateTicket(request)
	if err != nil {
		return
	}

	err = u.ticket.Deactivate(request, user, nonce)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
			u.clearCredential()
		}
		return
	}

	u.clearCredential()

	return nil
}
