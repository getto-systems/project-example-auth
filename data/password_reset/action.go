package password_reset

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	CreateSession interface {
		Create(request request.Request, user user.User, login user.Login) (Session, Destination, Token, error)
	}

	PushSendTokenJob interface {
		Push(request request.Request, session Session, dest Destination, token Token) error
	}

	SendToken interface {
		Send() error
	}

	GetStatus interface {
		Get(request request.Request, login user.Login, session Session) (Destination, Status, error)
	}

	Validate interface {
		Validate(request request.Request, login user.Login, token Token) (user.User, Session, ticket.Expiration, error)
	}

	CloseSession interface {
		Close(request request.Request, session Session) error
	}
)
