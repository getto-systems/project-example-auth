package password_reset

import (
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	SessionRepository interface {
		FindStatus(Session) (SessionData, Status, bool, error)
		FindSession(Token) (Session, SessionData, bool, error)
		CheckClosedSessionExists(Token) (bool, error)

		CreateSession(SessionGenerator, SessionData) (Session, Token, error)
		CloseSession(Session) error

		UpdateStatusToProcessing(Session) error
		UpdateStatusToFailed(Session, error) error
		UpdateStatusToSuccess(Session, Destination) error
	}

	DestinationRepository interface {
		FindDestination(user.User) (Destination, bool, error)

		RegisterDestination(user.User, Destination) error
	}
)
