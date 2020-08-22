package auth

import (
	"github.com/getto-systems/project-example-auth/request"
)

type (
	Renew struct {
		Usecase
	}
	RenewHandler interface {
		RenewRequest() (request.Request, error)
		RenewResponse(error)
	}
)

func NewRenew(u Usecase) Renew {
	return Renew{Usecase: u}
}

func (u Renew) Renew(handler RenewHandler) {
	handler.RenewResponse(u.renew(handler))
}
func (u Renew) renew(handler RenewHandler) (err error) {
	request, err := handler.RenewRequest()
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

	ticket, err := u.ticket.Extend(request, user, nonce)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	err = u.issueCredential(request, ticket)
	if err != nil {
		return
	}

	return nil
}
