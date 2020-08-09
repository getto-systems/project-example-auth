package infra

import (
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	SessionRepository interface {
		FindSessionDataAndDestinationAndStatus(password_reset.Session) (password_reset.SessionData, password_reset.Destination, password_reset.Status, bool, error)
		FindSession(password_reset.Token) (password_reset.Session, password_reset.SessionData, bool, error)
		CheckClosedSessionExists(password_reset.Token) (bool, error)

		CreateSession(SessionGenerator, password_reset.SessionData, password_reset.Destination) (password_reset.Session, password_reset.Token, error)
		CloseSession(password_reset.Session) error

		UpdateStatusToSending(password_reset.Session, request.RequestedAt) error
		UpdateStatusToFailed(password_reset.Session, request.RequestedAt, error) error
		UpdateStatusToComplete(password_reset.Session, request.RequestedAt) error
	}

	DestinationRepository interface {
		FindDestination(user.User) (password_reset.Destination, bool, error)

		RegisterDestination(user.User, password_reset.Destination) error
	}
)
