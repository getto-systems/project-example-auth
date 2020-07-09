package pubsub

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"
)

type SyncPubSub struct {
	authenticated []user.UserAuthenticatedEventHandler
	ticket        []user.UserTicketAuthEventHandler
	password      []user.UserPasswordAuthEventHandler
}

func NewSyncPubSub() *SyncPubSub {
	return &SyncPubSub{}
}

func (sub *SyncPubSub) SubscribeAuthenticated(handler user.UserAuthenticatedEventHandler) {
	sub.authenticated = append(sub.authenticated, handler)
}

func (sub *SyncPubSub) SubscribeTicketAuth(handler user.UserTicketAuthEventHandler) {
	sub.ticket = append(sub.ticket, handler)
}

func (sub *SyncPubSub) SubscribePasswordAuth(handler user.UserPasswordAuthEventHandler) {
	sub.password = append(sub.password, handler)
}

func (pub *SyncPubSub) Authenticated(request data.Request, ticket data.Ticket) {
	for _, handler := range pub.authenticated {
		handler.Authenticated(request, ticket)
	}
}

func (pub *SyncPubSub) TicketIssueFailed(request data.Request, ticket data.Ticket, err error) {
	for _, handler := range pub.authenticated {
		handler.TicketIssueFailed(request, ticket, err)
	}
}

func (pub *SyncPubSub) SignedTicketParsing(request data.Request) {
	for _, handler := range pub.ticket {
		handler.SignedTicketParsing(request)
	}
}

func (pub *SyncPubSub) SignedTicketParseFailed(request data.Request, err error) {
	for _, handler := range pub.ticket {
		handler.SignedTicketParseFailed(request, err)
	}
}

func (pub *SyncPubSub) PasswordMatching(request data.Request, userID data.UserID) {
	for _, handler := range pub.password {
		handler.PasswordMatching(request, userID)
	}
}

func (pub *SyncPubSub) PasswordMatchFailed(request data.Request, userID data.UserID, err error) {
	for _, handler := range pub.password {
		handler.PasswordMatchFailed(request, userID, err)
	}
}
