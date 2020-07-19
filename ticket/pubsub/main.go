package pubsub

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type PubSub struct {
	handlers []ticket.EventHandler
}

func NewPubSub() *PubSub {
	return &PubSub{}
}

func (pubsub *PubSub) pub() ticket.EventPublisher {
	return pubsub
}

func (pubsub *PubSub) Subscribe(handler ticket.EventHandler) {
	pubsub.handlers = append(pubsub.handlers, handler)
}

func (pubsub *PubSub) IssueApiToken(request data.Request, user data.User, roles data.Roles, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.IssueApiToken(request, user, roles, expires)
	}
}

func (pubsub *PubSub) IssueApiTokenFailed(request data.Request, user data.User, roles data.Roles, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueApiTokenFailed(request, user, roles, expires, err)
	}
}

func (pubsub *PubSub) IssueContentToken(request data.Request, user data.User, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.IssueContentToken(request, user, expires)
	}
}

func (pubsub *PubSub) IssueContentTokenFailed(request data.Request, user data.User, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueContentTokenFailed(request, user, expires, err)
	}
}

func (pubsub *PubSub) ExtendTicket(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.ExtendTicket(request, nonce, user, expires)
	}
}

func (pubsub *PubSub) ExtendTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.ExtendTicketFailed(request, nonce, user, expires, err)
	}
}

func (pubsub *PubSub) IssueTicket(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit) {
	for _, handler := range pubsub.handlers {
		handler.IssueTicket(request, user, expires, limit)
	}
}

func (pubsub *PubSub) IssueTicketFailed(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueTicketFailed(request, user, expires, limit, err)
	}
}

func (pubsub *PubSub) ShrinkTicket(request data.Request, nonce ticket.Nonce, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.ShrinkTicket(request, nonce, user)
	}
}

func (pubsub *PubSub) ShrinkTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.ShrinkTicketFailed(request, nonce, user, err)
	}
}

func (pubsub *PubSub) ValidateTicket(request data.Request) {
	for _, handler := range pubsub.handlers {
		handler.ValidateTicket(request)
	}
}

func (pubsub *PubSub) ValidateTicketFailed(request data.Request, err error) {
	for _, handler := range pubsub.handlers {
		handler.ValidateTicketFailed(request, err)
	}
}

func (pubsub *PubSub) AuthenticatedByTicket(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.AuthenticatedByTicket(request, user)
	}
}
