package password_reset

import (
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	SendTokenJobQueue interface {
		PushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination, password_reset.Token) error
		FetchSendTokenJob() (request.Request, password_reset.Session, password_reset.Destination, password_reset.Token, error)
	}
)
