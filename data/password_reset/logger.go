package password_reset

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Logger interface {
		CreateSessionLogger
		PushSendTokenJobLogger
		SendTokenLogger
		GetStatusLogger
		ValidateLogger
	}

	CreateSessionLogger interface {
		TryToCreateSession(request.Request, user.User, user.Login, time.Expires)
		FailedToCreateSession(request.Request, user.User, user.Login, time.Expires, error)
		CreateSession(request.Request, user.User, user.Login, time.Expires, Session, Destination)
	}

	PushSendTokenJobLogger interface {
		TryToPushSendTokenJob(request.Request, Session, Destination)
		FailedToPushSendTokenJob(request.Request, Session, Destination, error)
		PushSendTokenJob(request.Request, Session, Destination)
	}

	SendTokenLogger interface {
		TryToSendToken(request.Request, Session, Destination)
		FailedToSendToken(request.Request, Session, Destination, error)
		SendToken(request.Request, Session, Destination)
	}

	GetStatusLogger interface {
		TryToGetStatus(request.Request, Session)
		FailedToGetStatus(request.Request, Session, error)
		GetStatus(request.Request, Session, Status)
	}

	ValidateLogger interface {
		TryToValidateToken(request.Request, user.Login)
		FailedToValidateToken(request.Request, user.Login, error)
		FailedToValidateTokenBecauseForbidden(request.Request, user.Login, error)
		AuthByToken(request.Request, user.Login, user.User)
	}
)
