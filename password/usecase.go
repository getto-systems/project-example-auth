package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"

	"errors"
)

var (
	ErrValidateFailed = errors.New("password-validate-failed")
	ErrLoginNotFound  = errors.New("password-login-not-found")
	ErrRegisterFailed = errors.New("password-register-failed")
)

type Usecase interface {
	Validate(data.Request, Login, RawPassword) (
		ticket.Ticket,
		ticket.Nonce,
		ticket.ApiToken,
		ticket.ContentToken,
		data.Expires,
		error,
	)

	GetLogin(data.Request, ticket.Ticket, ticket.Nonce) (Login, error)
	Register(data.Request, ticket.Ticket, ticket.Nonce, Login, RegisterParam) error
}

type RegisterParam struct {
	OldPassword RawPassword
	NewPassword RawPassword
}

type (
	EventPublisher interface {
		RegisterEventPublisher
		ValidateEventPublisher
	}

	EventHandler interface {
		EventPublisher
	}

	DB interface {
		RegisterDB
		ValidateDB
	}
)
