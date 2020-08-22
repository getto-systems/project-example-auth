package password_reset_core

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (action action) Validate(request request.Request, login user.Login, token password_reset.Token) (_ user.User, _ password_reset.Session, _ credential.TicketExtendSecond, err error) {
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
			err = password_reset.ErrValidateAlreadyClosed
			action.logger.FailedToValidateTokenBecauseSessionClosed(request, login, err)
			return
		}

		err = password_reset.ErrValidateNotFoundSession
		action.logger.FailedToValidateTokenBecauseSessionNotFound(request, login, err)
		return
	}
	if data.Login().ID() != login.ID() {
		err = password_reset.ErrValidateMatchFailedLogin
		action.logger.FailedToValidateTokenBecauseLoginMatchFailed(request, login, err)
		return
	}
	if data.Expires().Expired(request) {
		err = password_reset.ErrValidateAlreadyExpired
		action.logger.FailedToValidateTokenBecauseSessionExpired(request, login, err)
		return
	}

	action.logger.AuthByToken(request, login, data.User())
	return data.User(), session, action.ticketExtendSecond, nil
}
