package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type (
	LoginID string
	Login   struct {
		id LoginID
	}

	RawPassword    string
	HashedPassword []byte

	PasswordGenerator interface {
		GeneratePassword(RawPassword) (HashedPassword, error)
	}

	PasswordMatcher interface {
		MatchPassword(HashedPassword, RawPassword) (bool, error)
	}

	PasswordEncrypter interface {
		PasswordGenerator
		PasswordMatcher
	}

	PasswordRepository interface {
		FindUser(Login) (data.User, bool, error)
		FindLogin(data.User) (Login, bool, error)
		FindPassword(Login) (data.User, HashedPassword, bool, error)

		RegisterPassword(data.User, HashedPassword) error
	}
)

func NewLogin(loginID LoginID) Login {
	return Login{
		id: loginID,
	}
}

func (login Login) ID() LoginID {
	return login.id
}
