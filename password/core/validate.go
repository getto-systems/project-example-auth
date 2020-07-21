package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type validator struct {
	pub  password.ValidateEventPublisher
	repo validateRepository
}

func newValidator(
	pub password.ValidateEventPublisher,
	db password.ValidateDB,
	matcher password.Matcher,
) validator {
	return validator{
		pub:  pub,
		repo: newValidateRepository(db, matcher),
	}
}

func (validator validator) validate(request data.Request, login password.Login, password password.RawPassword) (data.User, error) {
	validator.pub.ValidatePassword(request, login)

	err := checkPassword(password)
	if err != nil {
		validator.pub.ValidatePasswordFailed(request, login, err)
		return data.User{}, err
	}

	user, err := validator.repo.matchPassword(login, password)
	if err != nil {
		validator.pub.ValidatePasswordFailed(request, login, err)
		return data.User{}, err
	}

	validator.pub.AuthenticatedByPassword(request, login, user)

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

func (repo validateRepository) matchPassword(login password.Login, password password.RawPassword) (data.User, error) {
	user, hashed, err := repo.db.FindPasswordByLogin(login)
	if err != nil {
		return data.User{}, err
	}

	err = repo.matcher.MatchPassword(hashed, password)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}
