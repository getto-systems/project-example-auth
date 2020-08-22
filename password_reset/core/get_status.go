package password_reset_core

import (
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (action action) GetStatus(request request.Request, login user.Login, session password_reset.Session) (_ password_reset.Destination, _ password_reset.Status, err error) {
	action.logger.TryToGetStatus(request, session)

	data, dest, status, found, err := action.sessions.FindSessionDataAndDestinationAndStatus(session)
	if err != nil {
		action.logger.FailedToGetStatus(request, session, err)
		return
	}
	if !found {
		err = password_reset.ErrGetStatusNotFoundSession
		action.logger.FailedToGetStatusBecauseSessionNotFound(request, session, err)
		return
	}
	if data.Login().ID() != login.ID() {
		err = password_reset.ErrGetStatusMatchFailedLogin
		action.logger.FailedToGetStatusBecauseLoginMatchFailed(request, session, err)
		return
	}

	action.logger.GetStatus(request, session, status)
	return dest, status, nil
}
