package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

type Validater struct {
	pub    validateEventPublisher
	signer Signer
}

type validateEventPublisher interface {
	ValidateTicket(data.Request)
	ValidateTicketFailed(data.Request, error)
	AuthenticatedByTicket(data.Request, data.User)
}

func NewValidater(pub validateEventPublisher, signer Signer) Validater {
	return Validater{
		pub:    pub,
		signer: signer,
	}
}

func (validater Validater) validate(request data.Request, ticket Ticket, nonce Nonce) (data.User, error) {
	validater.pub.ValidateTicket(request)

	ticketNonce, user, expires, err := validater.signer.Parse(ticket)
	if err != nil {
		validater.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	if ticketNonce != nonce {
		err = errors.New("ticket nonce different")
		validater.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	if request.RequestedAt().Expired(expires) {
		err = errors.New("ticket already expired")
		validater.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	validater.pub.AuthenticatedByTicket(request, user)

	return user, nil
}
