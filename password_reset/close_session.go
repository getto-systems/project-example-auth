package password_reset

import (
	infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

var (
	errCloseNotFoundSession  = data.NewError("PasswordReset.Close", "NotFound.Session")
	errCloseMatchFailedLogin = data.NewError("PasswordReset.Close", "MatchFailed.Login")
	errCloseAlreadyExpired   = data.NewError("PasswordReset.Close", "AlreadyExpired")
)

type CloseSession struct {
	logger   infra.CloseSessionLogger
	sessions infra.SessionRepository
}

func NewCloseSession(logger infra.CloseSessionLogger, sessions infra.SessionRepository) CloseSession {
	return CloseSession{
		logger:   logger,
		sessions: sessions,
	}
}

func (action CloseSession) Close(request request.Request, session password_reset.Session) (err error) {
	action.logger.TryToCloseSession(request, session)

	err = action.sessions.CloseSession(session)
	if err != nil {
		action.logger.FailedToCloseSession(request, session, err)
		return
	}

	action.logger.CloseSession(request, session)
	return nil
}
