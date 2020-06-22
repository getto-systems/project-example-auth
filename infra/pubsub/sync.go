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

func (pub *SyncPubSub) Authenticated(request data.Request, ticket data.Ticket) {
	for _, handler := range pub.handlers {
		handler.Authenticated(request, ticket)
	}
}

func (pub *SyncPubSub) Authorized(request data.Request, ticket data.Ticket, resource data.Resource) {
	for _, handler := range pub.handlers {
		handler.Authorized(request, ticket, resource)
	}
}

func (pub *SyncPubSub) AuthorizeFailed(request data.Request, ticket data.Ticket, resource data.Resource, err error) {
	for _, handler := range pub.handlers {
		handler.AuthorizeFailed(request, ticket, resource, err)
	}
}

func (pub *SyncPubSub) TicketRenewing(request data.Request, ticket data.Ticket) {
	for _, handler := range pub.handlers {
		handler.TicketRenewing(request, ticket)
	}
}

func (pub *SyncPubSub) TicketRenewFailed(request data.Request, ticket data.Ticket, err error) {
	for _, handler := range pub.handlers {
		handler.TicketRenewFailed(request, ticket, err)
	}
}

func (pub *SyncPubSub) TicketRenewed(request data.Request, ticket data.Ticket) {
	for _, handler := range pub.handlers {
		handler.TicketRenewed(request, ticket)
	}
}

func (pub *SyncPubSub) PasswordMatching(request data.Request, user data.User) {
	for _, handler := range pub.handlers {
		handler.PasswordMatching(request, user)
	}
}

func (pub *SyncPubSub) PasswordMatchFailed(request data.Request, user data.User, err error) {
	for _, handler := range pub.handlers {
		handler.PasswordMatchFailed(request, user, err)
	}
}

func (pub *SyncPubSub) TicketIssueFailed(request data.Request, user data.User, err error) {
	for _, handler := range pub.handlers {
		handler.TicketIssueFailed(request, user, err)
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
