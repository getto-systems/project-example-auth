package user

import (
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type UserTicketAuth struct {
	pub  UserTicketAuthEventPublisher
	sign ticket.TicketSign

	request data.Request
}

func (user UserTicketAuth) Authenticate(signedTicket data.SignedTicket) (data.Ticket, error) {
	user.pub.SignedTicketParsing(user.request)

	ticket, err := user.sign.Parse(signedTicket)
	if err != nil {
		user.pub.SignedTicketParseFailed(user.request, err)
		return data.Ticket{}, err
	}

	if ticket.Expires.Expired(user.request.RequestedAt) {
		user.pub.SignedTicketParseFailed(user.request, ErrTicketAlreadyExpired)
		return data.Ticket{}, ErrTicketAlreadyExpired
	}

	return ticket, nil
}

type UserTicketAuthEventPublisher interface {
	SignedTicketParsing(data.Request)
	SignedTicketParseFailed(data.Request, error)
}

type UserTicketAuthEventHandler interface {
	UserTicketAuthEventPublisher
}

type UserTicketAuthFactory struct {
	pub  UserTicketAuthEventPublisher
	sign ticket.TicketSign
}

func NewUserTicketAuthFactory(pub UserTicketAuthEventPublisher, sign ticket.TicketSign) UserTicketAuthFactory {
	return UserTicketAuthFactory{
		pub:  pub,
		sign: sign,
	}
}

func (f UserTicketAuthFactory) New(request data.Request) UserTicketAuth {
	return UserTicketAuth{
		pub:  f.pub,
		sign: f.sign,

		request: request,
	}
}
