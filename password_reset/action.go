package password_reset

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	Action interface {
		CreateSession(request.Request, user.User, user.Login) (Session, Destination, Token, error)
		PushSendTokenJob(request.Request, Session, Destination, Token) error
		SendToken() error
		GetStatus(request.Request, user.Login, Session) (Destination, Status, error)
		Validate(request.Request, user.Login, Token) (user.User, Session, credential.TicketExtendSecond, error)
		CloseSession(request.Request, Session) error
	}
)
