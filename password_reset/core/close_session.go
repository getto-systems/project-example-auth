package password_reset_core

import (
	"github.com/getto-systems/project-example-id/misc/errors"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

var (
	errCloseNotFoundSession  = errors.NewError("PasswordReset.Close", "NotFound.Session")
	errCloseMatchFailedLogin = errors.NewError("PasswordReset.Close", "MatchFailed.Login")
	errCloseAlreadyExpired   = errors.NewError("PasswordReset.Close", "AlreadyExpired")
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
