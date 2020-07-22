package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type validator struct {
	logger  password.ValidateLogger
	matcher password.Matcher
	repo    validateRepository
}

func newValidator(
	logger password.ValidateLogger,
	db password.ValidateDB,
	matcher password.Matcher,
) validator {
	return validator{
		logger:  logger,
		matcher: matcher,
		repo:    newValidateRepository(db),
	}
}

func (validator validator) validate(
	request data.Request,
	login password.Login,
	raw password.RawPassword,
) (user data.User, err error) {
	validator.logger.TryToValidate(request, login)
	defer func() {
		if err == nil {
			validator.logger.AuthedByPassword(request, login, user)
		} else {
			validator.logger.FailedToValidate(request, login, err)
		}
	}()

	err = checkPassword(raw)
	if err != nil {
		return
	}

	pass, err := validator.repo.findPassword(login)
	if err != nil {
		return
	}

	return pass.Match(validator.matcher, raw)
}

type validateRepository struct {
	db password.ValidateDB
}

func newValidateRepository(db password.ValidateDB) validateRepository {
	return validateRepository{
		db: db,
	}
}

func (repo validateRepository) findPassword(login password.Login) (pass password.Password, err error) {
	passwordSlice, err := repo.db.FilterPassword(login)
	if err != nil {
		return
	}

	if len(passwordSlice) == 0 {
		err = password.ErrPasswordNotFound
		return
	}

	return passwordSlice[0], nil
}
