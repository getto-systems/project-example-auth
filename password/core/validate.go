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

func (validator validator) validate(request data.Request, user data.User, password password.RawPassword) error {
	validator.pub.ValidatePassword(request, user)

	err := checkPassword(password)
	if err != nil {
		validator.pub.ValidatePasswordFailed(request, user, err)
		return err
	}

	err = validator.repo.matchPassword(user, password)
	if err != nil {
		validator.pub.ValidatePasswordFailed(request, user, err)
		return err
	}

	validator.pub.AuthenticatedByPassword(request, user)

	return nil
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

func (repo validateRepository) matchPassword(user data.User, password password.RawPassword) error {
	hashed, err := repo.db.FindUserPassword(user)
	if err != nil {
		return err
	}

	return repo.matcher.MatchPassword(hashed, password)
}
