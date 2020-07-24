package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type validator struct {
	logger    password.ValidateLogger
	passwords password.PasswordRepository
	matcher   password.PasswordMatcher
}

func newValidator(
	logger password.ValidateLogger,
	passwords password.PasswordRepository,
	matcher password.PasswordMatcher,
) validator {
	return validator{
		logger:    logger,
		passwords: passwords,
		matcher:   matcher,
	}
}

func (validator validator) validate(
	request data.Request,
	login password.Login,
	raw password.RawPassword,
) (_ data.User, err error) {
	validator.logger.TryToValidate(request, login)
	defer func() {
		if err != nil {
			validator.logger.FailedToValidate(request, login, err)
		}
	}()

	err = check(raw)
	if err != nil {
		return
	}

	user, hashed, found, err := validator.passwords.FindPassword(login)
	if err != nil {
		return
	}
	if !found {
		err = errPasswordNotFoundPassword
		return
	}

	matched, err := validator.matcher.MatchPassword(hashed, raw)
	if err != nil {
		return
	}
	if !matched {
		err = errPasswordNotMatched
		return
	}

	validator.logger.AuthedByPassword(request, login, user)
	return user, nil
}
