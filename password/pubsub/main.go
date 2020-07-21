package pubsub

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type PubSub struct {
	handlers []password.EventHandler
}

func NewPubSub() *PubSub {
	return &PubSub{}
}

func (pubsub *PubSub) pub() password.EventPublisher {
	return pubsub
}

func (pubsub *PubSub) Subscribe(handler password.EventHandler) {
	pubsub.handlers = append(pubsub.handlers, handler)
}

func (pubsub *PubSub) GetLogin(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.GetLogin(request, user)
	}
}

func (pubsub *PubSub) LoginNotFound(request data.Request, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.LoginNotFound(request, user, err)
	}
}

func (pubsub *PubSub) RegisterPassword(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.RegisterPassword(request, user)
	}
}

func (pubsub *PubSub) RegisterPasswordFailed(request data.Request, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.RegisterPasswordFailed(request, user, err)
	}
}

func (pubsub *PubSub) PasswordRegistered(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.PasswordRegistered(request, user)
	}
}

func (pubsub *PubSub) ValidatePassword(request data.Request, login password.Login) {
	for _, handler := range pubsub.handlers {
		handler.ValidatePassword(request, login)
	}
}

func (pubsub *PubSub) ValidatePasswordFailed(request data.Request, login password.Login, err error) {
	for _, handler := range pubsub.handlers {
		handler.ValidatePasswordFailed(request, login, err)
	}
}

func (pubsub *PubSub) AuthenticatedByPassword(request data.Request, login password.Login, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.AuthenticatedByPassword(request, login, user)
	}
}
