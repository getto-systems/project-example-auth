package password_reset

import (
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	SessionRepository interface {
		FindStatus(Session) (SessionData, Status, bool, error)
		FindSession(Token) (SessionData, bool, error)

		RegisterSession(SessionGenerator, SessionData) (Session, Token, error)

		UpdateStatusToProcessing(Session) error
		UpdateStatusToFailed(Session, error) error
		UpdateStatusToSuccess(Session, Destination) error
	}

	DestinationRepository interface {
		FindDestination(user.User) (Destination, bool, error)

		RegisterDestination(user.User, Destination) error
	}
)
