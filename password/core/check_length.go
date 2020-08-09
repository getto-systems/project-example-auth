package password_core

import (
	"github.com/getto-systems/project-example-id/misc/errors"

	"github.com/getto-systems/project-example-id/password"
)

const (
	PASSWORD_MAX_BYTES = 72 // bcrypt max length
)

var (
	errCheckLengthEmpty   = errors.NewError("Password.Check", "Length.Empty")
	errCheckLengthTooLong = errors.NewError("Password.Check", "Length.TooLong")
)

func checkLength(raw password.RawPassword) (err error) {
	if len(raw) == 0 {
		err = errCheckLengthEmpty
		return
	}

	if len(raw) > PASSWORD_MAX_BYTES {
		err = errCheckLengthTooLong
		return
	}

	return nil
}
