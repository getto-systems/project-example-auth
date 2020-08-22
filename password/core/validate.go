package password_core

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (action action) Validate(request request.Request, user user.User, raw password.RawPassword) (_ credential.TicketExtendSecond, err error) {
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
		err = password.ErrValidateNotFoundPassword
		action.logger.FailedToValidateBecausePasswordNotFound(request, user, err)
		return
	}

	matched, err := action.matcher.MatchPassword(hashed, raw)
	if err != nil {
		action.logger.FailedToValidate(request, user, err)
		return
	}
	if !matched {
		err = password.ErrValidateMatchFailed
		action.logger.FailedToValidateBecausePasswordMatchFailed(request, user, err)
		return
	}

	action.logger.AuthByPassword(request, user)
	return action.extendSecond, nil
}
