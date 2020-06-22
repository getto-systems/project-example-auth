package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type AuthenticatedUser struct {
	pub UserEventPublisher

	ticket data.Ticket
}

func (user AuthenticatedUser) Ticket() data.Ticket {
	return user.ticket
}

type User struct {
	pub UserEventPublisher

	user data.User
}

func (user User) Authenticated(ticket data.Ticket) AuthenticatedUser {
	return AuthenticatedUser{
		pub:    user.pub,
		ticket: ticket,
	}
}

func (user User) UserID() data.UserID {
	return user.user.UserID
}

type UnauthorizedUser struct {
	pub UserEventPublisher
}

func (user UnauthorizedUser) Authenticated(ticket data.Ticket) AuthenticatedUser {
	return AuthenticatedUser{
		pub:    user.pub,
		ticket: ticket,
	}
}

func (user AuthenticatedUser) Authenticated(request data.Request) {
	user.pub.Authenticated(request, user.ticket)
}

func (user AuthenticatedUser) Authorized(request data.Request, resource data.Resource) {
	user.pub.Authorized(request, user.ticket, resource)
}

func (user AuthenticatedUser) AuthorizeFailed(request data.Request, resource data.Resource, err error) {
	user.pub.AuthorizeFailed(request, user.ticket, resource, err)
}

func (user AuthenticatedUser) TicketRenewing(request data.Request) {
	user.pub.TicketRenewing(request, user.ticket)
}

func (user AuthenticatedUser) TicketRenewFailed(request data.Request, err error) {
	user.pub.TicketRenewFailed(request, user.ticket, err)
}

func (user AuthenticatedUser) TicketRenewed(request data.Request) {
	user.pub.TicketRenewed(request, user.ticket)
}

func (user User) PasswordMatching(request data.Request) {
	user.pub.PasswordMatching(request, user.user)
}

func (user User) PasswordMatchFailed(request data.Request, err error) {
	user.pub.PasswordMatchFailed(request, user.user, err)
}

func (user User) TicketIssueFailed(request data.Request, err error) {
	user.pub.TicketIssueFailed(request, user.user, err)
}

func (user UnauthorizedUser) Authorizing(request data.Request, resource data.Resource) {
	user.pub.Authorizing(request, resource)
}

func (user UnauthorizedUser) AuthorizeTokenParseFailed(request data.Request, resource data.Resource, err error) {
	user.pub.AuthorizeTokenParseFailed(request, resource, err)
}

type UserEventPublisher interface {
	Authenticated(data.Request, data.Ticket)
	Authorized(data.Request, data.Ticket, data.Resource)
	AuthorizeFailed(data.Request, data.Ticket, data.Resource, error)

	TicketRenewing(data.Request, data.Ticket)
	TicketRenewFailed(data.Request, data.Ticket, error)
	TicketRenewed(data.Request, data.Ticket)

	PasswordMatching(data.Request, data.User)
	PasswordMatchFailed(data.Request, data.User, error)

	TicketIssueFailed(data.Request, data.User, error)

	Authorizing(data.Request, data.Resource)
	AuthorizeTokenParseFailed(data.Request, data.Resource, error)
}

type UserEventHandler interface {
	UserEventPublisher
}

type UserEventSubscriber interface {
	Subscribe(UserEventHandler)
}

type UserFactory struct {
	pub UserEventPublisher
}

func (f UserFactory) Authenticated(ticket data.Ticket) AuthenticatedUser {
	return AuthenticatedUser{
		pub:    f.pub,
		ticket: ticket,
	}
}

func (f UserFactory) New(userID data.UserID) User {
	return User{
		pub: f.pub,
		user: data.User{
			UserID: userID,
		},
	}
}

func (f UserFactory) Unauthorized() UnauthorizedUser {
	return UnauthorizedUser{
		pub: f.pub,
	}
}

func NewUserFactory(pub UserEventPublisher) UserFactory {
	return UserFactory{
		pub: pub,
	}
}
