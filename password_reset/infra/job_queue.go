package infra

import (
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
)

type (
	SendTokenJobQueue interface {
		PushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination, password_reset.Token) error
		FetchSendTokenJob() (request.Request, password_reset.Session, password_reset.Destination, password_reset.Token, error)
	}
)
