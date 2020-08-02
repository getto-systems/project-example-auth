package password

type (
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
)
