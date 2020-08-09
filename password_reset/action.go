package password_reset

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Action interface {
		CreateSession(request request.Request, user user.User, login user.Login) (Session, Destination, Token, error)
		PushSendTokenJob(request request.Request, session Session, dest Destination, token Token) error
		SendToken() error
		GetStatus(request request.Request, login user.Login, session Session) (Destination, Status, error)
		Validate(request request.Request, login user.Login, token Token) (user.User, Session, credential.Expiration, error)
		CloseSession(request request.Request, session Session) error
	}
)
