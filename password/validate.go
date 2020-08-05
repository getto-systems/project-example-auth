package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateNotFoundPassword = data.NewError("Password.Validate", "NotFound.Password")
	errValidateNotMatched       = data.NewError("Password.Validate", "NotMatched")
)

type Validate struct {
	logger    password.ValidateLogger
	matcher   password.PasswordMatcher
	passwords password.PasswordRepository
}

func NewValidate(logger password.ValidateLogger, matcher password.PasswordMatcher, passwords password.PasswordRepository) Validate {
	return Validate{
		logger:    logger,
		passwords: passwords,
		matcher:   matcher,
	}
}

func (action Validate) Validate(request request.Request, user user.User, raw password.RawPassword) (err error) {
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
		action.logger.FailedToValidate(request, user, err)
		return
	}

	matched, err := action.matcher.MatchPassword(hashed, raw)
	if err != nil {
		action.logger.FailedToValidate(request, user, err)
		return
	}
	if !matched {
		err = errValidateNotMatched
		action.logger.FailedToValidate(request, user, err)
		return
	}

	action.logger.AuthByPassword(request, user)
	return nil
}
