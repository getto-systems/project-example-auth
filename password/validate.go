package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateNotFoundPassword = data.NewError("Password.Validate", "NotFound.Password")
	errValidateMatchFailed      = data.NewError("Password.Validate", "MatchFailed")
)

type Validate struct {
	logger    password.ValidateLogger
	exp       ticket.Expiration
	matcher   password.PasswordMatcher
	passwords password.PasswordRepository
}

func NewValidate(logger password.ValidateLogger, exp ticket.Expiration, matcher password.PasswordMatcher, passwords password.PasswordRepository) Validate {
	return Validate{
		logger:    logger,
		exp:       exp,
		matcher:   matcher,
		passwords: passwords,
	}
}

func (action Validate) Validate(request request.Request, user user.User, raw password.RawPassword) (_ ticket.Expiration, err error) {
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
