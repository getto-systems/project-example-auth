package password_reset

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		CreateSession(request.Request, user.User, user.Login) (Session, Destination, Token, error)
		PushSendTokenJob(request.Request, Session, Destination, Token) error
		SendToken() error
		GetStatus(request.Request, user.Login, Session) (Destination, Status, error)
		Validate(request.Request, user.Login, Token) (user.User, Session, expiration.ExtendSecond, error)
		CloseSession(request.Request, Session) error
	}
)
