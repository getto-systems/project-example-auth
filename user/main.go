package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type User struct {
	pub UserEventPublisher

	userID data.UserID
}

func (user User) UserID() data.UserID {
	return user.userID
}

func (user User) Authenticated(request data.Request) {
	user.pub.Authenticated(request, user.userID)
}

func (user User) TicketRenewing(request data.Request) {
	user.pub.TicketRenewing(request, user.userID)
}

func (user User) TicketRenewFailed(request data.Request, err error) {
	user.pub.TicketRenewFailed(request, user.userID, err)
}

func (user User) TicketRenewed(request data.Request) {
	user.pub.TicketRenewed(request, user.userID)
}

func (user User) PasswordMatching(request data.Request) {
	user.pub.PasswordMatching(request, user.userID)
}

func (user User) PasswordMatchFailed(request data.Request, err error) {
	user.pub.PasswordMatchFailed(request, user.userID, err)
}

func (user User) TicketIssueFailed(request data.Request, err error) {
	user.pub.TicketIssueFailed(request, user.userID, err)
}

func (user User) AuthorizeFailed(request data.Request, resource data.Resource, err error) {
	user.pub.AuthorizeFailed(request, user.userID, resource, err)
}

func (user User) Authorized(request data.Request, resource data.Resource) {
	user.pub.Authorized(request, user.userID, resource)
}

type UnauthorizedUser struct {
	pub UserEventPublisher
}

func (user UnauthorizedUser) Authorizing(request data.Request, resource data.Resource) {
	user.pub.Authorizing(request, resource)
}

func (user UnauthorizedUser) AuthorizeTokenParseFailed(request data.Request, resource data.Resource, err error) {
	user.pub.AuthorizeTokenParseFailed(request, resource, err)
}

type UserEventPublisher interface {
	Authenticated(data.Request, data.UserID)
	Authorized(data.Request, data.UserID, data.Resource)

	TicketRenewing(data.Request, data.UserID)
	TicketRenewFailed(data.Request, data.UserID, error)
	TicketRenewed(data.Request, data.UserID)

	PasswordMatching(data.Request, data.UserID)
	PasswordMatchFailed(data.Request, data.UserID, error)
	TicketIssueFailed(data.Request, data.UserID, error)

	Authorizing(data.Request, data.Resource)
	AuthorizeTokenParseFailed(data.Request, data.Resource, error)
	AuthorizeFailed(data.Request, data.UserID, data.Resource, error)
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

func (f UserFactory) New(userID data.UserID) User {
	return User{
		pub:    f.pub,
		userID: userID,
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
