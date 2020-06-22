package pubsub

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"
)

type SyncPubSub struct {
	handlers []user.UserEventHandler
}

func NewSyncPubSub() *SyncPubSub {
	return &SyncPubSub{}
}

func (sub *SyncPubSub) Subscribe(handler user.UserEventHandler) {
	sub.handlers = append(sub.handlers, handler)
}

func (pub *SyncPubSub) Authenticated(request data.Request, userID data.UserID, profile data.Profile) {
	for _, handler := range pub.handlers {
		handler.Authenticated(request, userID, profile)
	}
}

func (pub *SyncPubSub) Authorized(request data.Request, userID data.UserID, profile data.Profile, resource data.Resource) {
	for _, handler := range pub.handlers {
		handler.Authorized(request, userID, profile, resource)
	}
}

func (pub *SyncPubSub) TicketRenewing(request data.Request, userID data.UserID) {
	for _, handler := range pub.handlers {
		handler.TicketRenewing(request, userID)
	}
}

func (pub *SyncPubSub) TicketRenewFailed(request data.Request, userID data.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.TicketRenewFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) TicketRenewed(request data.Request, userID data.UserID) {
	for _, handler := range pub.handlers {
		handler.TicketRenewed(request, userID)
	}
}

func (pub *SyncPubSub) PasswordMatching(request data.Request, userID data.UserID) {
	for _, handler := range pub.handlers {
		handler.PasswordMatching(request, userID)
	}
}

func (pub *SyncPubSub) PasswordMatchFailed(request data.Request, userID data.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.PasswordMatchFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) TicketIssueFailed(request data.Request, userID data.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.TicketIssueFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) Authorizing(request data.Request, resource data.Resource) {
	for _, handler := range pub.handlers {
		handler.Authorizing(request, resource)
	}
}

func (pub *SyncPubSub) AuthorizeTokenParseFailed(request data.Request, resource data.Resource, err error) {
	for _, handler := range pub.handlers {
		handler.AuthorizeTokenParseFailed(request, resource, err)
	}
}

func (pub *SyncPubSub) AuthorizeFailed(request data.Request, userID data.UserID, resource data.Resource, err error) {
	for _, handler := range pub.handlers {
		handler.AuthorizeFailed(request, userID, resource, err)
	}
}
