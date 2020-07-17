package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type Validater struct {
	pub  validateEventPublisher
	repo validateRepository
}

type validateEventPublisher interface {
	ValidatePassword(data.Request, data.User)
	ValidatePasswordFailed(data.Request, data.User, error)
	AuthenticatedByPassword(data.Request, data.User)
}

type validateDB interface {
	FindUserPassword(data.User) (data.HashedPassword, error)
}

func NewValidater(
	pub validateEventPublisher,
	db validateDB,
	matcher Matcher,
) Validater {
	return Validater{
		pub:  pub,
		repo: newValidateRepository(db, matcher),
	}
}

func (validater Validater) validate(request data.Request, user data.User, password data.RawPassword) error {
	validater.pub.ValidatePassword(request, user)

	err := checkPassword(password)
	if err != nil {
		validater.pub.ValidatePasswordFailed(request, user, err)
		return err
	}

	err = validater.repo.matchPassword(user, password)
	if err != nil {
		validater.pub.ValidatePasswordFailed(request, user, err)
		return err
	}

	validater.pub.AuthenticatedByPassword(request, user)

	return nil
}

type validateRepository struct {
	db      validateDB
	matcher Matcher
}

func newValidateRepository(db validateDB, matcher Matcher) validateRepository {
	return validateRepository{
		db:      db,
		matcher: matcher,
	}
}

type Matcher interface {
	MatchPassword(data.HashedPassword, data.RawPassword) error
}

func (repo validateRepository) matchPassword(user data.User, password data.RawPassword) error {
	hashed, err := repo.db.FindUserPassword(user)
	if err != nil {
		return err
	}

	return repo.matcher.MatchPassword(hashed, password)
}
