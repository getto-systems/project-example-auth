package password_reset_core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/password_reset"
)

var (
	errCloseNotFoundSession  = data.NewError("PasswordReset.Close", "NotFound.Session")
	errCloseMatchFailedLogin = data.NewError("PasswordReset.Close", "MatchFailed.Login")
	errCloseAlreadyExpired   = data.NewError("PasswordReset.Close", "AlreadyExpired")
)

func (action action) CloseSession(request request.Request, session password_reset.Session) (err error) {
	action.logger.TryToCloseSession(request, session)

	err = action.sessions.CloseSession(session)
	if err != nil {
		action.logger.FailedToCloseSession(request, session, err)
		return
	}

	action.logger.CloseSession(request, session)
	return nil
}
