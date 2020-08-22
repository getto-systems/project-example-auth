package password

import (
	"errors"
)

var (
	ErrCheckLengthEmpty   = errors.New("Password.Check/Length.Empty")
	ErrCheckLengthTooLong = errors.New("Password.Check/Length.TooLong")

	ErrValidateNotFoundPassword = errors.New("Password.Validate/NotFound.Password")
	ErrValidateMatchFailed      = errors.New("Password.Validate/MatchFailed")
)

type (
	RawPassword    string
	HashedPassword []byte

	ChangeParam struct {
		OldPassword RawPassword
		NewPassword RawPassword
	}
)
