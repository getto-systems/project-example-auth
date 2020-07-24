package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type (
	Usecase interface {
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

		CreateResetSession(data.Request, Login) (ResetSession, error)
		GetResetStatus(data.Request, ResetSession) (ResetStatus, error)
		Reset(data.Request, Login, ResetToken, RawPassword) (
			ticket.Ticket,
			ticket.Nonce,
			ticket.ApiToken,
			ticket.ContentToken,
			data.Expires,
			error,
		)
	}

	RegisterParam struct {
		OldPassword RawPassword
		NewPassword RawPassword
	}

	Logger interface {
		ValidateLogger
		RegisterLogger
		ResetLogger
	}
)
