package password

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
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

	err := validater.checkPassword(password)
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
func (validater Validater) checkPassword(password data.RawPassword) error {
	if len(password) == 0 {
		return errors.New("password is empty")
	}
	if len([]byte(password)) > PASSWORD_BYTES_LIMIT {
		return errors.New("password is too long")
	}
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
