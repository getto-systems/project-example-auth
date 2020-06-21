package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrTicketRenewFailed = errors.New("ticket renew failed")
)

type RenewAuthenticator struct {
	issuerRepository IssuerRepository

	user    user.User
	request data.Request
}

func (authenticator RenewAuthenticator) RenewTicket(ticket data.Ticket) (data.Token, error) {
	authenticator.ticketRenewing()

	issuer := authenticator.issuerRepository.New(ticket)

	token, err := issuer.Renew(authenticator.request.RequestedAt)
	if err != nil {
		authenticator.ticketRenewFailed(err)
		return nil, ErrTicketRenewFailed
	}

	authenticator.ticketRenewed()

	return token, nil
}

func (authenticator RenewAuthenticator) ticketRenewing() {
	authenticator.user.TicketRenewing(authenticator.request)
}

func (authenticator RenewAuthenticator) ticketRenewFailed(err error) {
	authenticator.user.TicketRenewFailed(authenticator.request, err)
}

func (authenticator RenewAuthenticator) ticketRenewed() {
	authenticator.user.TicketRenewed(authenticator.request)
}

type RenewAuthenticatorFactory struct {
	issuerRepository IssuerRepository
	userFactory      user.UserFactory
}

func (f RenewAuthenticatorFactory) New(userID data.UserID, request data.Request) RenewAuthenticator {
	return RenewAuthenticator{
		issuerRepository: f.issuerRepository,

		user:    f.userFactory.New(userID),
		request: request,
	}
}

func NewRenewAuthenticatorFactory(issuerRepository IssuerRepository, userFactory user.UserFactory) RenewAuthenticatorFactory {
	return RenewAuthenticatorFactory{
		issuerRepository: issuerRepository,
		userFactory:      userFactory,
	}
}
