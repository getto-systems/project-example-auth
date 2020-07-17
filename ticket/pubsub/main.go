package pubsub

import (
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type TicketPubSub struct {
	handlers []ticket.EventHandler
}

func NewTicketPubSub() *TicketPubSub {
	return &TicketPubSub{}
}

func (pubsub *TicketPubSub) Publisher() ticket.EventPublisher {
	return pubsub
}

func (pubsub *TicketPubSub) Subscribe(handler ticket.EventHandler) {
	pubsub.handlers = append(pubsub.handlers, handler)
}

func (pubsub *TicketPubSub) IssueApiToken(request data.Request, user data.User, roles data.Roles, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.IssueApiToken(request, user, roles, expires)
	}
}

func (pubsub *TicketPubSub) IssueApiTokenFailed(request data.Request, user data.User, roles data.Roles, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueApiTokenFailed(request, user, roles, expires, err)
	}
}

func (pubsub *TicketPubSub) IssueContentToken(request data.Request, user data.User, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.IssueContentToken(request, user, expires)
	}
}

func (pubsub *TicketPubSub) IssueContentTokenFailed(request data.Request, user data.User, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueContentTokenFailed(request, user, expires, err)
	}
}

func (pubsub *TicketPubSub) ExtendTicket(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.ExtendTicket(request, nonce, user, expires)
	}
}

func (pubsub *TicketPubSub) ExtendTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.ExtendTicketFailed(request, nonce, user, expires, err)
	}
}

func (pubsub *TicketPubSub) IssueTicket(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit) {
	for _, handler := range pubsub.handlers {
		handler.IssueTicket(request, user, expires, limit)
	}
}

func (pubsub *TicketPubSub) IssueTicketFailed(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueTicketFailed(request, user, expires, limit, err)
	}
}

func (pubsub *TicketPubSub) ShrinkTicket(request data.Request, nonce ticket.Nonce, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.ShrinkTicket(request, nonce, user)
	}
}

func (pubsub *TicketPubSub) ShrinkTicketFailed(request data.Request, nonce ticket.Nonce, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.ShrinkTicketFailed(request, nonce, user, err)
	}
}

func (pubsub *TicketPubSub) ValidateTicket(request data.Request) {
	for _, handler := range pubsub.handlers {
		handler.ValidateTicket(request)
	}
}

func (pubsub *TicketPubSub) ValidateTicketFailed(request data.Request, err error) {
	for _, handler := range pubsub.handlers {
		handler.ValidateTicketFailed(request, err)
	}
}

func (pubsub *TicketPubSub) AuthenticatedByTicket(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.AuthenticatedByTicket(request, user)
	}
}
