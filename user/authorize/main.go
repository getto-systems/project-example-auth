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
	userFactory      user.UserFactory

	unauthorized user.UnauthorizedUser
	user         user.User

	token   data.Token
	ticket  data.Ticket
	request data.Request
}

func (authorizer Authorizer) IsAccessible(resource data.Resource) (data.Ticket, error) {
	authorizer.authorizing(resource)

	err := authorizer.parseToken()
	if err != nil {
		authorizer.authorizeTokenParseFailed(resource, err)
		return data.Ticket{}, ErrAuthorizeTokenParseFailed
	}

	err = authorizer.hasEnoughPermission(resource)
	if err != nil {
		authorizer.authorizeFailed(resource, err)
		return data.Ticket{}, ErrAuthorizeFailed
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

func (authorizer Authorizer) hasEnoughPermission(resource data.Resource) error {
	return authorizer.ticketAuthorizer.HasEnoughPermission(authorizer.ticket, authorizer.request, resource)
}

func (authorizer Authorizer) authorizing(resource data.Resource) {
	authorizer.unauthorized.Authorizing(authorizer.request, resource)
}

func (authorizer Authorizer) authorizeTokenParseFailed(resource data.Resource, err error) {
	authorizer.unauthorized.AuthorizeTokenParseFailed(authorizer.request, resource, err)
}

func (authorizer Authorizer) authorizeFailed(resource data.Resource, err error) {
	authorizer.user.AuthorizeFailed(authorizer.request, resource, err)
}

func (authorizer Authorizer) authorized(resource data.Resource) {
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

func (f AuthorizerFactory) New(token data.Token, request data.Request) Authorizer {
	return Authorizer{
		ticketAuthorizer: f.ticketAuthorizer,
		userFactory:      f.userFactory,

		unauthorized: f.userFactory.Unauthorized(),

		token:   token,
		request: request,
	}
}
