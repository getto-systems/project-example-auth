package password

import (
	"github.com/getto-systems/project-example-id/z_external/errors"
)

var (
	ErrCheck = errors.NewCategory("PasswordCheck")

	ErrCheckLengthEmpty   = errors.NewErrorAsCategory("Password.Check", "Length.Empty", ErrCheck)
	ErrCheckLengthTooLong = errors.NewErrorAsCategory("Password.Check", "Length.TooLong", ErrCheck)

	ErrValidateNotFoundPassword = errors.NewError("Password.Validate", "NotFound.Password")
	ErrValidateMatchFailed      = errors.NewError("Password.Validate", "MatchFailed")
)

type (
	RawPassword    string
	HashedPassword []byte

	ChangeParam struct {
		OldPassword RawPassword
		NewPassword RawPassword
	}
)
