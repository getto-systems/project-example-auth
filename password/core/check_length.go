package password_core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

const (
	PASSWORD_MAX_BYTES = 72 // bcrypt max length
)

var (
	errCheckLengthEmpty   = data.NewError("Password.Check", "Length.Empty")
	errCheckLengthTooLong = data.NewError("Password.Check", "Length.TooLong")
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
