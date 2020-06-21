package user

import (
	"github.com/getto-systems/project-example-id/basic"
)

type User struct {
	pub UserEventPublisher

	userID basic.UserID
}

func (user User) UserID() basic.UserID {
	return user.userID
}

func (user User) Authenticated(request basic.Request) {
	user.pub.Authenticated(request, user.userID)
}

func (user User) TicketRenewing(request basic.Request) {
	user.pub.TicketRenewing(request, user.userID)
}

func (user User) TicketRenewFailed(request basic.Request, err error) {
	user.pub.TicketRenewFailed(request, user.userID, err)
}

func (user User) TicketRenewed(request basic.Request) {
	user.pub.TicketRenewed(request, user.userID)
}

func (user User) PasswordMatching(request basic.Request) {
	user.pub.PasswordMatching(request, user.userID)
}

func (user User) PasswordMatchFailed(request basic.Request, err error) {
	user.pub.PasswordMatchFailed(request, user.userID, err)
}

func (user User) TicketIssueFailed(request basic.Request, err error) {
	user.pub.TicketIssueFailed(request, user.userID, err)
}

func (user User) AuthorizeFailed(request basic.Request, resource basic.Resource, err error) {
	user.pub.AuthorizeFailed(request, user.userID, resource, err)
}

func (user User) Authorized(request basic.Request, resource basic.Resource) {
	user.pub.Authorized(request, user.userID, resource)
}

type UnauthorizedUser struct {
	pub UserEventPublisher
}

func (user UnauthorizedUser) Authorizing(request basic.Request, resource basic.Resource) {
	user.pub.Authorizing(request, resource)
}

func (user UnauthorizedUser) AuthorizeTokenParseFailed(request basic.Request, resource basic.Resource, err error) {
	user.pub.AuthorizeTokenParseFailed(request, resource, err)
}

type UserEventPublisher interface {
	Authenticated(basic.Request, basic.UserID)
	Authorized(basic.Request, basic.UserID, basic.Resource)

	TicketRenewing(basic.Request, basic.UserID)
	TicketRenewFailed(basic.Request, basic.UserID, error)
	TicketRenewed(basic.Request, basic.UserID)

	PasswordMatching(basic.Request, basic.UserID)
	PasswordMatchFailed(basic.Request, basic.UserID, error)
	TicketIssueFailed(basic.Request, basic.UserID, error)

	Authorizing(basic.Request, basic.Resource)
	AuthorizeTokenParseFailed(basic.Request, basic.Resource, error)
	AuthorizeFailed(basic.Request, basic.UserID, basic.Resource, error)
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

func (f UserFactory) New(userID basic.UserID) User {
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
