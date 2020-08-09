package password_reset_core

import (
	"github.com/getto-systems/project-example-id/errors"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errCreateSessionNotFoundDestination = errors.NewError("PasswordReset.CreateSession", "NotFound.Destination")
)

func (action action) CreateSession(request request.Request, user user.User, login user.Login) (_ password_reset.Session, _ password_reset.Destination, _ password_reset.Token, err error) {
	expires := action.exp.Expires(request)

	action.logger.TryToCreateSession(request, user, login, expires)

	dest, found, err := action.destinations.FindDestination(user)
	if err != nil {
		action.logger.FailedToCreateSession(request, user, login, expires, err)
		return
	}
	if !found {
		err = errCreateSessionNotFoundDestination
		action.logger.FailedToCreateSessionBecauseDestinationNotFound(request, user, login, expires, err)
		return
	}

	session, token, err := action.sessions.CreateSession(
		action.sessionGenerator,
		password_reset.NewSessionData(user, login, request.RequestedAt(), expires),
		dest,
	)
	if err != nil {
		action.logger.FailedToCreateSession(request, user, login, expires, err)
		return
	}

	action.logger.CreateSession(request, user, login, expires, session, dest)
	return session, dest, token, nil
}
