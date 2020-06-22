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

	authenticatedUser user.AuthenticatedUser

	request data.Request
}

func (authenticator RenewAuthenticator) RenewTicket(ticket data.Ticket) (data.Token, error) {
	authenticator.fireTicketRenewing()

	issuer := authenticator.issuerRepository.New(ticket)

	token, err := issuer.Renew(authenticator.request.RequestedAt)
	if err != nil {
		authenticator.fireTicketRenewFailed(err)
		return nil, ErrTicketRenewFailed
	}

	authenticator.fireTicketRenewed()

	return token, nil
}

func (authenticator RenewAuthenticator) fireTicketRenewing() {
	authenticator.authenticatedUser.TicketRenewing(authenticator.request)
}

func (authenticator RenewAuthenticator) fireTicketRenewFailed(err error) {
	authenticator.authenticatedUser.TicketRenewFailed(authenticator.request, err)
}

func (authenticator RenewAuthenticator) fireTicketRenewed() {
	authenticator.authenticatedUser.TicketRenewed(authenticator.request)
}

type RenewAuthenticatorFactory struct {
	issuerRepository IssuerRepository
	userFactory      user.UserFactory
}

func (f RenewAuthenticatorFactory) New(ticket data.Ticket, request data.Request) RenewAuthenticator {
	return RenewAuthenticator{
		issuerRepository: f.issuerRepository,

		authenticatedUser: f.userFactory.Authenticated(ticket),
		request:           request,
	}
}

func NewRenewAuthenticatorFactory(issuerRepository IssuerRepository, userFactory user.UserFactory) RenewAuthenticatorFactory {
	return RenewAuthenticatorFactory{
		issuerRepository: issuerRepository,
		userFactory:      userFactory,
	}
}
