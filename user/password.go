package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type Password interface {
	Match(data.RawPassword) error
}

type HashedPassword struct {
	matcher PasswordMatcher
	hashed  data.HashedPassword
}

func (p HashedPassword) Match(password data.RawPassword) error {
	return p.matcher.MatchPassword(p.hashed, password)
}

type NullPassword struct {
}

func (NullPassword) Match(password data.RawPassword) error {
	return ErrUserPasswordNotFound
}
