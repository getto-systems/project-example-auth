package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"

	"errors"
)

var (
	ErrValidateFailed = errors.New("password-validate-failed")
	ErrRegisterFailed = errors.New("password-register-failed")
)

type Usecase interface {
	Validate(data.Request, data.User, RawPassword) (
		ticket.Ticket,
		ticket.Nonce,
		ticket.ApiToken,
		ticket.ContentToken,
		data.Expires,
		error,
	)
	Register(data.Request, ticket.Ticket, ticket.Nonce, RegisterParam) error
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
