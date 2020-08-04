package password_reset

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errCreateSessionNotFoundDestination = errors.NewError("PasswordReset.CreateSession", "NotFound.Destination")
)

type CreateSession struct {
	logger       password_reset.CreateSessionLogger
	exp          expiration
	gen          password_reset.SessionGenerator
	sessions     password_reset.SessionRepository
	destinations password_reset.DestinationRepository
}

func NewCreateSession(logger password_reset.CreateSessionLogger, exp time.Second, gen password_reset.SessionGenerator, sessions password_reset.SessionRepository, destinations password_reset.DestinationRepository) CreateSession {
	return CreateSession{
		logger:       logger,
		exp:          newExpiration(exp),
		gen:          gen,
		sessions:     sessions,
		destinations: destinations,
	}
}

func (action CreateSession) Create(request request.Request, user user.User, login user.Login) (_ password_reset.Session, _ password_reset.Destination, _ password_reset.Token, err error) {
	requestedAt := request.RequestedAt()
	expires := action.exp.Expires(requestedAt)

	action.logger.TryToCreateSession(request, user, login, expires)

	dest, found, err := action.destinations.FindDestination(user)
	if err != nil {
		action.logger.FailedToCreateSession(request, user, login, expires, err)
		return
	}
	if !found {
		err = errCreateSessionNotFoundDestination
		action.logger.FailedToCreateSession(request, user, login, expires, err)
		return
	}

	session, token, err := action.sessions.RegisterSession(
		action.gen,
		password_reset.NewSessionData(user, login, requestedAt, expires),
	)
	if err != nil {
		action.logger.FailedToCreateSession(request, user, login, expires, err)
		return
	}

	action.logger.CreateSession(request, user, login, expires, session, dest)
	return session, dest, token, nil
}