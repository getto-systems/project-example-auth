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

func (pubsub *PubSub) GetLogin(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.GetLogin(request, user)
	}
}

func (pubsub *PubSub) GetLoginFailed(request data.Request, user data.User, err error) {
	for _, handler := range pubsub.handlers {
		handler.GetLoginFailed(request, user, err)
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

func (pubsub *PubSub) RegisteredPassword(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.RegisteredPassword(request, user)
	}
}

func (pubsub *PubSub) IssueResetToken(request data.Request, login password.Login, expires data.Expires) {
	for _, handler := range pubsub.handlers {
		handler.IssueResetToken(request, login, expires)
	}
}

func (pubsub *PubSub) IssueResetTokenFailed(request data.Request, login password.Login, expires data.Expires, err error) {
	for _, handler := range pubsub.handlers {
		handler.IssueResetTokenFailed(request, login, expires, err)
	}
}

func (pubsub *PubSub) GetResetStatus(request data.Request, reset password.Reset) {
	for _, handler := range pubsub.handlers {
		handler.GetResetStatus(request, reset)
	}
}

func (pubsub *PubSub) GetResetStatusFailed(request data.Request, reset password.Reset, err error) {
	for _, handler := range pubsub.handlers {
		handler.GetResetStatusFailed(request, reset, err)
	}
}

func (pubsub *PubSub) ValidateResetToken(request data.Request) {
	for _, handler := range pubsub.handlers {
		handler.ValidateResetToken(request)
	}
}

func (pubsub *PubSub) ValidateResetTokenFailed(request data.Request, err error) {
	for _, handler := range pubsub.handlers {
		handler.ValidateResetTokenFailed(request, err)
	}
}

func (pubsub *PubSub) AuthenticatedByResetToken(request data.Request, user data.User) {
	for _, handler := range pubsub.handlers {
		handler.AuthenticatedByResetToken(request, user)
	}
}
