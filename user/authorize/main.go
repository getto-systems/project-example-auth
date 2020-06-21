package authorize

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/basic"

	"errors"
)

var (
	ErrAuthorizeTokenParseFailed = errors.New("token parse failed")
	ErrAuthorizeFailed           = errors.New("authorize failed")
)

type Authorizer struct {
	ticketAuthorizer user.TicketAuthorizer
	userFactory      user.UserFactory

	unauthorized user.UnauthorizedUser
	user         user.User

	token   basic.Token
	ticket  basic.Ticket
	request basic.Request
}

func (authorizer Authorizer) IsAccessible(resource basic.Resource) (basic.Ticket, error) {
	authorizer.authorizing(resource)

	err := authorizer.parseToken()
	if err != nil {
		authorizer.authorizeTokenParseFailed(resource, err)
		return basic.Ticket{}, ErrAuthorizeTokenParseFailed
	}

	err = authorizer.hasEnoughPermission(resource)
	if err != nil {
		authorizer.authorizeFailed(resource, err)
		return basic.Ticket{}, ErrAuthorizeFailed
	}

	authorizer.authorized(resource)

	return authorizer.ticket, nil
}

func (authorizer Authorizer) parseToken() error {
	ticket, err := authorizer.ticketAuthorizer.DecodeToken(authorizer.token)
	if err != nil {
		return ErrAuthorizeTokenParseFailed
	}

	authorizer.ticket = ticket
	authorizer.user = authorizer.userFactory.New(ticket.Profile.UserID)

	return nil
}

func (authorizer Authorizer) hasEnoughPermission(resource basic.Resource) error {
	return authorizer.ticketAuthorizer.HasEnoughPermission(authorizer.ticket, authorizer.request, resource)
}

func (authorizer Authorizer) authorizing(resource basic.Resource) {
	authorizer.unauthorized.Authorizing(authorizer.request, resource)
}

func (authorizer Authorizer) authorizeTokenParseFailed(resource basic.Resource, err error) {
	authorizer.unauthorized.AuthorizeTokenParseFailed(authorizer.request, resource, err)
}

func (authorizer Authorizer) authorizeFailed(resource basic.Resource, err error) {
	authorizer.user.AuthorizeFailed(authorizer.request, resource, err)
}

func (authorizer Authorizer) authorized(resource basic.Resource) {
	authorizer.user.Authorized(authorizer.request, resource)
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

func (f AuthorizerFactory) New(token basic.Token, request basic.Request) Authorizer {
	return Authorizer{
		ticketAuthorizer: f.ticketAuthorizer,
		userFactory:      f.userFactory,

		unauthorized: f.userFactory.Unauthorized(),

		token:   token,
		request: request,
	}
}
