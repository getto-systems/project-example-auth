package password_reset

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	SessionRepository interface {
		FindSessionDataAndDestinationAndStatus(Session) (SessionData, Destination, Status, bool, error)
		FindSession(Token) (Session, SessionData, bool, error)
		CheckClosedSessionExists(Token) (bool, error)

		CreateSession(SessionGenerator, SessionData, Destination) (Session, Token, error)
		CloseSession(Session) error

		UpdateStatusToSending(Session, time.RequestedAt) error
		UpdateStatusToFailed(Session, time.RequestedAt, error) error
		UpdateStatusToComplete(Session, time.RequestedAt) error
	}

	DestinationRepository interface {
		FindDestination(user.User) (Destination, bool, error)

		RegisterDestination(user.User, Destination) error
	}
)
