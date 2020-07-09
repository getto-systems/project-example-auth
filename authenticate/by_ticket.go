package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"
)

type AuthByTicket struct {
	authenticated user.UserAuthenticatedFactory
	ticketAuth    user.UserTicketAuthFactory
}

func (auth AuthByTicket) Authenticate(request data.Request, signedTicket data.SignedTicket) (data.Ticket, data.SignedTicket, error) {
	ticketAuthUser := auth.ticketAuth.New(request)
	ticket, err := ticketAuthUser.Authenticate(signedTicket)
	if err != nil {
		return data.Ticket{}, nil, ErrTicketAuthFailed
	}

	authenticatedUser := auth.authenticated.New(request, ticket.Profile.UserID)
	return authenticatedUser.IssueTicket()
}

func NewAuthByTicket(authenticated user.UserAuthenticatedFactory, ticketAuth user.UserTicketAuthFactory) AuthByTicket {
	return AuthByTicket{
		authenticated: authenticated,
		ticketAuth:    ticketAuth,
	}
}
