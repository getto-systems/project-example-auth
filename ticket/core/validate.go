package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"

	"errors"
)

type validator struct {
	pub    ticket.ValidateEventPublisher
	signer ticket.Signer
}

func newValidator(
	pub ticket.ValidateEventPublisher,
	signer ticket.Signer,
) validator {
	return validator{
		pub:    pub,
		signer: signer,
	}
}

func (validator validator) validate(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce) (data.User, error) {
	validator.pub.ValidateTicket(request)

	ticketNonce, user, expires, err := validator.signer.Parse(ticket)
	if err != nil {
		validator.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	if ticketNonce != nonce {
		err = errors.New("ticket nonce different")
		validator.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	if request.RequestedAt().Expired(expires) {
		err = errors.New("ticket already expired")
		validator.pub.ValidateTicketFailed(request, err)
		return data.User{}, err
	}

	validator.pub.AuthenticatedByTicket(request, user)

	return user, nil
}
