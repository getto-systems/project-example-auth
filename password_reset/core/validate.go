package password_reset_core

import (
	"github.com/getto-systems/project-example-id/_misc/errors"
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errValidateNotFoundSession  = errors.NewError("PasswordReset.Validate", "NotFound.Session")
	errValidateMatchFailedLogin = errors.NewError("PasswordReset.Validate", "MatchFailed.Login")
	errValidateAlreadyExpired   = errors.NewError("PasswordReset.Validate", "AlreadyExpired")
	errValidateAlreadyClosed    = errors.NewError("PasswordReset.Validate", "AlreadyClosed")
)

func (action action) Validate(request request.Request, login user.Login, token password_reset.Token) (_ user.User, _ password_reset.Session, _ expiration.ExtendSecond, err error) {
	action.logger.TryToValidateToken(request, login)

	session, data, found, err := action.sessions.FindSession(token)
	if err != nil {
		action.logger.FailedToValidateToken(request, login, err)
		return
	}
	if !found {
		found, err = action.sessions.CheckClosedSessionExists(token)
		if err != nil {
			action.logger.FailedToValidateToken(request, login, err)
			return
		}
		if found {
			err = errValidateAlreadyClosed
			action.logger.FailedToValidateTokenBecauseSessionClosed(request, login, err)
			return
		}

		err = errValidateNotFoundSession
		action.logger.FailedToValidateTokenBecauseSessionNotFound(request, login, err)
		return
	}
	if data.Login().ID() != login.ID() {
		err = errValidateMatchFailedLogin
		action.logger.FailedToValidateTokenBecauseLoginMatchFailed(request, login, err)
		return
	}
	if request.Expired(data.Expires()) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateTokenBecauseSessionExpired(request, login, err)
		return
	}

	action.logger.AuthByToken(request, login, data.User())
	return data.User(), session, action.extendSecond, nil
}