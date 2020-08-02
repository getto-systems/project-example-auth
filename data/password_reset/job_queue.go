package password_reset

import (
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	SendTokenJobQueue interface {
		PushSendTokenJob(request.Request, Session, Destination, Token) error
		FetchSendTokenJob() (request.Request, Session, Destination, Token, error)
	}
)
