package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/basic"
)

type PasswordRepository struct {
	db PasswordDB

	passwordMatcherFactory user.PasswordMatcherFactory
}

type PasswordDB interface {
	UserPassword(basic.UserID) (basic.HashedPassword, error)
}

func (repo PasswordRepository) Find(userID basic.UserID) user.PasswordMatcher {
	password, err := repo.db.UserPassword(userID)
	if err != nil {
		return repo.passwordMatcherFactory.NotFound(err)
	}

	return repo.passwordMatcherFactory.New(password)
}

func NewPasswordRepository(db PasswordDB, passwordMatcherFactory user.PasswordMatcherFactory) PasswordRepository {
	return PasswordRepository{
		db: db,

		passwordMatcherFactory: passwordMatcherFactory,
	}
}
