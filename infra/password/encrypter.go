package password

import (
	"github.com/getto-systems/project-example-id/data/password"
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
