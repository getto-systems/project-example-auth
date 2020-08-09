package _usecase

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
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
	credential, err := u.renew(handler)
	u.handleCredential(credential, err)
	handler.RenewResponse(err)
}
func (u Renew) renew(handler RenewHandler) (_ credential.Credential, err error) {
	ticket, err := u.getTicket()
	if err != nil {
		return
	}

	request, err := handler.RenewRequest()
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

	expires, err := u.ticket.Extend(request, user, ticket)
	if err != nil {
		return
	}

	return u.issueCredential(request, user, ticket.Nonce(), expires)
}
