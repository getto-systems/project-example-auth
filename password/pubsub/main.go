package pubsub

import (
	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/data"
)

type PasswordPubSub struct {
	handlers []password.EventHandler
}

func NewPasswordPubSub() *PasswordPubSub {
	return &PasswordPubSub{}
}

func (pubsub *PasswordPubSub) Publisher() password.EventPublisher {
	return pubsub
}

func (pubsub *PasswordPubSub) Subscribe(handler password.EventHandler) {
	pubsub.handlers = append(pubsub.handlers, handler)
}

func (pubsub *PasswordPubSub) ValidatePassword(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.ValidatePassword(request, user)
	}
}

func (pubsub *PasswordPubSub) ValidatePasswordFailed(request data.Request, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.ValidatePasswordFailed(request, user, err)
	}
}

func (pubsub *PasswordPubSub) PasswordRegistered(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.PasswordRegistered(request, user)
	}
}

func (pubsub *PasswordPubSub) VerifyPassword(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.VerifyPassword(request, user)
	}
}

func (pubsub *PasswordPubSub) VerifyPasswordFailed(request data.Request, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.VerifyPasswordFailed(request, user, err)
	}
}

func (pubsub *PasswordPubSub) AuthenticatedByPassword(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.AuthenticatedByPassword(request, user)
	}
}
