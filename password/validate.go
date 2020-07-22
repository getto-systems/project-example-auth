package password

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrPasswordNotFound = errors.New("password not found")
	ErrLoginNotFound    = errors.New("login not found")
)

type ValidateLogger interface {
	TryToValidate(data.Request, Login)
	FailedToValidate(data.Request, Login, error)
	AuthedByPassword(data.Request, Login, data.User)
}

type ValidateDB interface {
	FilterPassword(Login) ([]Password, error)
}

type Matcher interface {
	MatchPassword(HashedPassword, RawPassword) error
}

func (password Password) Match(matcher Matcher, raw RawPassword) error {
	return matcher.MatchPassword(password.hashed, raw)
}
