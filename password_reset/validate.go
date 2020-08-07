package password_reset

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateNotFoundSession  = data.NewError("PasswordReset.Validate", "NotFound.Session")
	errValidateMatchFailedLogin = data.NewError("PasswordReset.Validate", "MatchFailed.Login")
	errValidateAlreadyExpired   = data.NewError("PasswordReset.Validate", "AlreadyExpired")
	errValidateAlreadyClosed    = data.NewError("PasswordReset.Validate", "AlreadyClosed")
)

type Validate struct {
	logger   password_reset.ValidateLogger
	exp      ticket.Expiration
	sessions password_reset.SessionRepository
}

func NewValidate(logger password_reset.ValidateLogger, exp ticket.Expiration, sessions password_reset.SessionRepository) Validate {
	return Validate{
		logger:   logger,
		exp:      exp,
		sessions: sessions,
	}
}

func (action Validate) Validate(request request.Request, login user.Login, token password_reset.Token) (_ user.User, _ password_reset.Session, _ ticket.Expiration, err error) {
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
	if request.RequestedAt().Expired(data.Expires()) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateTokenBecauseSessionExpired(request, login, err)
		return
	}

	action.logger.AuthByToken(request, login, data.User())
	return data.User(), session, action.exp, nil
}
