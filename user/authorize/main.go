package authorize

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrAuthorizeTokenParseFailed = errors.New("token parse failed")
	ErrAuthorizeFailed           = errors.New("authorize failed")
)

type Authorizer struct {
	ticketAuthorizer user.TicketAuthorizer

	user    user.UnauthorizedUser
	request data.Request
}

func (authorizer *Authorizer) HasEnoughPermission(token data.Token, resource data.Resource) (data.Ticket, error) {
	authorizer.user.Authorizing(authorizer.request, resource)

	ticket, err := authorizer.parseToken(token)
	if err != nil {
		authorizer.user.AuthorizeTokenParseFailed(authorizer.request, resource, err)
		return data.Ticket{}, ErrAuthorizeTokenParseFailed
	}

	user := authorizer.user.Authenticated(ticket)

	err = authorizer.hasEnoughPermission(ticket, resource)
	if err != nil {
		user.AuthorizeFailed(authorizer.request, resource, err)
		return data.Ticket{}, ErrAuthorizeFailed
	}

	user.Authorized(authorizer.request, resource)

	return ticket, nil
}

func (authorizer *Authorizer) parseToken(token data.Token) (data.Ticket, error) {
	ticket, err := authorizer.ticketAuthorizer.DecodeToken(token)
	if err != nil {
		return data.Ticket{}, ErrAuthorizeTokenParseFailed
	}

	return ticket, nil
}

func (authorizer *Authorizer) hasEnoughPermission(ticket data.Ticket, resource data.Resource) error {
	return authorizer.ticketAuthorizer.HasEnoughPermission(ticket, authorizer.request, resource)
}

type AuthorizerFactory struct {
	ticketAuthorizer user.TicketAuthorizer
	userFactory      user.UserFactory
}

func NewAuthorizerFactory(ticketAuthorizer user.TicketAuthorizer, userFactory user.UserFactory) AuthorizerFactory {
	return AuthorizerFactory{
		ticketAuthorizer: ticketAuthorizer,
		userFactory:      userFactory,
	}
}

func (f AuthorizerFactory) New(request data.Request) Authorizer {
	return Authorizer{
		ticketAuthorizer: f.ticketAuthorizer,

		user:    f.userFactory.Unauthorized(),
		request: request,
	}
}
