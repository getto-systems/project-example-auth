package password_reset

import (
	infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errGetStatusNotFoundSession  = data.NewError("PasswordReset.GetStatus", "NotFound.Session")
	errGetStatusMatchFailedLogin = data.NewError("PasswordReset.GetStatus", "MatchFailed.Login")
)

type GetStatus struct {
	logger   infra.GetStatusLogger
	sessions infra.SessionRepository
}

func NewGetStatus(logger infra.GetStatusLogger, sessions infra.SessionRepository) GetStatus {
	return GetStatus{
		logger:   logger,
		sessions: sessions,
	}
}

func (action GetStatus) Get(request request.Request, login user.Login, session password_reset.Session) (_ password_reset.Destination, _ password_reset.Status, err error) {
	action.logger.TryToGetStatus(request, session)

	data, dest, status, found, err := action.sessions.FindSessionDataAndDestinationAndStatus(session)
	if err != nil {
		action.logger.FailedToGetStatus(request, session, err)
		return
	}
	if !found {
		err = errGetStatusNotFoundSession
		action.logger.FailedToGetStatusBecauseSessionNotFound(request, session, err)
		return
	}
	if data.Login().ID() != login.ID() {
		err = errGetStatusMatchFailedLogin
		action.logger.FailedToGetStatusBecauseLoginMatchFailed(request, session, err)
		return
	}

	action.logger.GetStatus(request, session, status)
	return dest, status, nil
}
