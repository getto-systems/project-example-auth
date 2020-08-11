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
	nonce, signature, err := u.getTicketNonceAndSignature()
	if err != nil {
		return
	}

	request, err := handler.RenewRequest()
	if err != nil {
		return
	}

	user, err := u.credential.ParseTicketSignature(request, nonce, signature)
	if err != nil {
		return
	}

	err = u.ticket.Validate(request, user, nonce)
	if err != nil {
		return
	}

	ticket, err := u.ticket.Extend(request, user, nonce)
	if err != nil {
		return
	}

	return u.issueCredential(request, ticket)
}
