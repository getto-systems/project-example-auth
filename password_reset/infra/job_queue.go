package infra

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	SendTokenJobQueue interface {
		PushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination, password_reset.Token) error
		FetchSendTokenJob() (request.Request, password_reset.Session, password_reset.Destination, password_reset.Token, error)
	}
)
