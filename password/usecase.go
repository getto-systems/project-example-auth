package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
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

	IssueReset(data.Request, Login) (Reset, error)
	GetResetStatus(data.Request, Reset) (ResetStatus, error)
	Reset(data.Request, Login, ResetToken, RawPassword) (
		ticket.Ticket,
		ticket.Nonce,
		ticket.ApiToken,
		ticket.ContentToken,
		data.Expires,
		error,
	)
}

type RegisterParam struct {
	OldPassword RawPassword
	NewPassword RawPassword
}

type (
	Logger interface {
		ValidateLogger
		RegisterLogger
		ResetLogger
	}

	DB interface {
		ValidateDB
		RegisterDB
		ResetDB
	}
)
