package password_reset

import (
	"github.com/getto-systems/project-example-id/data/password_reset"
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
		CloseSessionLogger
	}

	CreateSessionLogger interface {
		TryToCreateSession(request.Request, user.User, user.Login, time.Expires)
		FailedToCreateSession(request.Request, user.User, user.Login, time.Expires, error)
		FailedToCreateSessionBecauseDestinationNotFound(request.Request, user.User, user.Login, time.Expires, error)
		CreateSession(request.Request, user.User, user.Login, time.Expires, password_reset.Session, password_reset.Destination)
	}

	PushSendTokenJobLogger interface {
		TryToPushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination)
		FailedToPushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination, error)
		PushSendTokenJob(request.Request, password_reset.Session, password_reset.Destination)
	}

	SendTokenLogger interface {
		TryToSendToken(request.Request, password_reset.Session, password_reset.Destination)
		FailedToSendToken(request.Request, password_reset.Session, password_reset.Destination, error)
		SendToken(request.Request, password_reset.Session, password_reset.Destination)
	}

	GetStatusLogger interface {
		TryToGetStatus(request.Request, password_reset.Session)
		FailedToGetStatus(request.Request, password_reset.Session, error)
		FailedToGetStatusBecauseSessionNotFound(request.Request, password_reset.Session, error)
		FailedToGetStatusBecauseLoginMatchFailed(request.Request, password_reset.Session, error)
		GetStatus(request.Request, password_reset.Session, password_reset.Status)
	}

	ValidateLogger interface {
		TryToValidateToken(request.Request, user.Login)
		FailedToValidateToken(request.Request, user.Login, error)
		FailedToValidateTokenBecauseSessionNotFound(request.Request, user.Login, error)
		FailedToValidateTokenBecauseSessionClosed(request.Request, user.Login, error)
		FailedToValidateTokenBecauseSessionExpired(request.Request, user.Login, error)
		FailedToValidateTokenBecauseLoginMatchFailed(request.Request, user.Login, error)
		AuthByToken(request.Request, user.Login, user.User)
	}

	CloseSessionLogger interface {
		TryToCloseSession(request.Request, password_reset.Session)
		FailedToCloseSession(request.Request, password_reset.Session, error)
		CloseSession(request.Request, password_reset.Session)
	}
)
