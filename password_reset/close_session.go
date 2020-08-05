package password_reset

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

var (
	errCloseNotFoundSession = data.NewError("PasswordReset.Close", "NotFound.Session")
	errCloseNotMatchedLogin = data.NewError("PasswordReset.Close", "NotMatched.Login")
	errCloseAlreadyExpired  = data.NewError("PasswordReset.Close", "AlreadyExpired")
)

type CloseSession struct {
	logger   password_reset.CloseSessionLogger
	sessions password_reset.SessionRepository
}

func NewCloseSession(logger password_reset.CloseSessionLogger, sessions password_reset.SessionRepository) CloseSession {
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
