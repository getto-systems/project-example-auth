package password

import (
	"github.com/getto-systems/project-example-id/data"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
)

var (
	ErrPasswordEmpty      = newError("Password/Empty")
	ErrPasswordTooLong    = newError("Password/TooLong")
	ErrPasswordNotMatched = newError("Password/NotMatched")

	ErrPasswordNotFoundUser     = newError("Password/NotFound/User")
	ErrPasswordNotFoundLogin    = newError("Password/NotFound/Login")
	ErrPasswordNotFoundPassword = newError("Password/NotFound/Password")
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

func (raw RawPassword) Check() error {
	if len(raw) == 0 {
		return ErrPasswordEmpty
	}
	if len([]byte(raw)) > PASSWORD_BYTES_LIMIT {
		return ErrPasswordTooLong
	}
	return nil
}
