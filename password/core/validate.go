package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type validator struct {
	logger password.ValidateLogger
	repo   validateRepository
}

func newValidator(
	logger password.ValidateLogger,
	db password.ValidateDB,
	matcher password.Matcher,
) validator {
	return validator{
		logger: logger,
		repo:   newValidateRepository(db, matcher),
	}
}

func (validator validator) validate(request data.Request, login password.Login, password password.RawPassword) (data.User, error) {
	validator.logger.TryToValidate(request, login)

	user, err := validator.repo.match(login, password)
	if err != nil {
		validator.logger.FailedToValidate(request, login, err)
		return data.User{}, err
	}

	validator.logger.AuthedByPassword(request, login, user)

	return user, nil
}

type validateRepository struct {
	db      password.ValidateDB
	matcher password.Matcher
}

func newValidateRepository(db password.ValidateDB, matcher password.Matcher) validateRepository {
	return validateRepository{
		db:      db,
		matcher: matcher,
	}
}

func (repo validateRepository) match(login password.Login, raw password.RawPassword) (data.User, error) {
	err := checkPassword(raw)
	if err != nil {
		return data.User{}, err
	}

	password, err := repo.findPassword(login)
	if err != nil {
		return data.User{}, err
	}

	err = password.Match(repo.matcher, raw)
	if err != nil {
		return data.User{}, err
	}

	return password.User(), nil
}

func (repo validateRepository) findPassword(login password.Login) (password.Password, error) {
	passwordSlice, err := repo.db.FilterPassword(login)
	if err != nil {
		return password.Password{}, err
	}

	if len(passwordSlice) == 0 {
		return password.Password{}, password.ErrPasswordNotFound
	}

	return passwordSlice[0], nil
}
