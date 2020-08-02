package password_reset

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateNotFoundSession = errors.NewError("PasswordReset.Validate", "NotFound.Session")
	errValidateNotMatchedLogin = errors.NewError("PasswordReset.Validate", "NotMatched.Login")
	errValidateAlreadyExpired  = errors.NewError("PasswordReset.Validate", "AlreadyExpired")
)

type Validate struct {
	logger   password_reset.ValidateLogger
	sessions password_reset.SessionRepository
}

func NewValidate(logger password_reset.ValidateLogger, sessions password_reset.SessionRepository) Validate {
	return Validate{
		logger:   logger,
		sessions: sessions,
	}
}

func (action Validate) Validate(request request.Request, login user.Login, token password_reset.Token) (_ user.User, err error) {
	action.logger.TryToValidateToken(request, login)

	data, found, err := action.sessions.FindSession(token)
	if err != nil {
		action.logger.FailedToValidateToken(request, login, err)
		return
	}
	if !found {
		err = errValidateNotFoundSession
		action.logger.FailedToValidateToken(request, login, err)
		return
	}
	if data.Login().ID() != login.ID() {
		err = errValidateNotMatchedLogin
		action.logger.FailedToValidateTokenBecauseForbidden(request, login, err)
		return
	}
	if request.RequestedAt().Expired(data.Expires()) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateTokenBecauseForbidden(request, login, err)
		return
	}

	action.logger.AuthByToken(request, login, data.User())
	return data.User(), nil
}
