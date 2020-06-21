package pubsub

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/basic"
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

func (pub *SyncPubSub) Authenticated(request basic.Request, userID basic.UserID) {
	for _, handler := range pub.handlers {
		handler.Authenticated(request, userID)
	}
}

func (pub *SyncPubSub) Authorized(request basic.Request, userID basic.UserID, resource basic.Resource) {
	for _, handler := range pub.handlers {
		handler.Authorized(request, userID, resource)
	}
}

func (pub *SyncPubSub) TicketRenewing(request basic.Request, userID basic.UserID) {
	for _, handler := range pub.handlers {
		handler.TicketRenewing(request, userID)
	}
}

func (pub *SyncPubSub) TicketRenewFailed(request basic.Request, userID basic.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.TicketRenewFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) TicketRenewed(request basic.Request, userID basic.UserID) {
	for _, handler := range pub.handlers {
		handler.TicketRenewed(request, userID)
	}
}

func (pub *SyncPubSub) PasswordMatching(request basic.Request, userID basic.UserID) {
	for _, handler := range pub.handlers {
		handler.PasswordMatching(request, userID)
	}
}

func (pub *SyncPubSub) PasswordMatchFailed(request basic.Request, userID basic.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.PasswordMatchFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) TicketIssueFailed(request basic.Request, userID basic.UserID, err error) {
	for _, handler := range pub.handlers {
		handler.TicketIssueFailed(request, userID, err)
	}
}

func (pub *SyncPubSub) Authorizing(request basic.Request, resource basic.Resource) {
	for _, handler := range pub.handlers {
		handler.Authorizing(request, resource)
	}
}

func (pub *SyncPubSub) AuthorizeTokenParseFailed(request basic.Request, resource basic.Resource, err error) {
	for _, handler := range pub.handlers {
		handler.AuthorizeTokenParseFailed(request, resource, err)
	}
}

func (pub *SyncPubSub) AuthorizeFailed(request basic.Request, userID basic.UserID, resource basic.Resource, err error) {
	for _, handler := range pub.handlers {
		handler.AuthorizeFailed(request, userID, resource, err)
	}
}
