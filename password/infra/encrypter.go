package infra

import (
	"github.com/getto-systems/project-example-id/password"
)

type (
	PasswordGenerator interface {
		GeneratePassword(password.RawPassword) (password.HashedPassword, error)
	}

	PasswordMatcher interface {
		MatchPassword(password.HashedPassword, password.RawPassword) (bool, error)
	}

	PasswordEncrypter interface {
		PasswordGenerator
		PasswordMatcher
	}
)
