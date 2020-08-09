package password_core

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/errors"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errValidateNotFoundPassword = errors.NewError("Password.Validate", "NotFound.Password")
	errValidateMatchFailed      = errors.NewError("Password.Validate", "MatchFailed")
)

func (action action) Validate(request request.Request, user user.User, raw password.RawPassword) (_ ticket.Expiration, err error) {
	action.logger.TryToValidate(request, user)

	err = checkLength(raw)
	if err != nil {
		action.logger.FailedToValidateBecausePasswordCheckFailed(request, user, err)
		return
	}

	hashed, found, err := action.passwords.FindPassword(user)
	if err != nil {
		action.logger.FailedToValidate(request, user, err)
		return
	}
	if !found {
		err = errValidateNotFoundPassword
		action.logger.FailedToValidateBecausePasswordNotFound(request, user, err)
		return
	}

	matched, err := action.matcher.MatchPassword(hashed, raw)
	if err != nil {
		action.logger.FailedToValidate(request, user, err)
		return
	}
	if !matched {
		err = errValidateMatchFailed
		action.logger.FailedToValidateBecausePasswordMatchFailed(request, user, err)
		return
	}

	action.logger.AuthByPassword(request, user)
	return action.exp, nil
}
