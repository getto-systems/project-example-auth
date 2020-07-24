package password

import (
	"github.com/getto-systems/project-example-id/data"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
)

var (
	ErrPasswordEmpty   = newError("Password/Empty")
	ErrPasswordTooLong = newError("Password/TooLong")

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
		MatchPassword(HashedPassword, RawPassword) error
	}

	PasswordEncrypter interface {
		PasswordGenerator
		PasswordMatcher
	}

	PasswordRepository interface {
		FindUser(Login) (data.User, error)                     // 見つからない場合は ErrPasswordNotFoundUser
		FindLogin(data.User) (Login, error)                    // 見つからない場合は ErrPasswordNotFoundLogin
		FindPassword(Login) (data.User, HashedPassword, error) // 見つからない場合は ErrPasswordNotFoundPassword

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
