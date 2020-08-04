package password_reset

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errGetStatusNotFoundSession = errors.NewError("PasswordReset.GetStatus", "NotFound.Session")
	errGetStatusDifferentLogin  = errors.NewError("PasswordReset.GetStatus", "DifferentLogin")
)

type GetStatus struct {
	logger   password_reset.GetStatusLogger
	sessions password_reset.SessionRepository
}

func NewGetStatus(logger password_reset.GetStatusLogger, sessions password_reset.SessionRepository) GetStatus {
	return GetStatus{
		logger:   logger,
		sessions: sessions,
	}
}

func (action GetStatus) Get(request request.Request, login user.Login, session password_reset.Session) (_ password_reset.Status, err error) {
	action.logger.TryToGetStatus(request, session)

	data, status, found, err := action.sessions.FindStatus(session)
	if err != nil {
		action.logger.FailedToGetStatus(request, session, err)
		return
	}
	if !found {
		err = errGetStatusNotFoundSession
		return
	}
	if data.Login().ID() != login.ID() {
		err = errGetStatusDifferentLogin
		return
	}

	action.logger.GetStatus(request, session, status)
	return status, nil
}
