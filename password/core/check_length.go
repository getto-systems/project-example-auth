package password_core

import (
	"github.com/getto-systems/project-example-auth/password"
)

const (
	PASSWORD_MAX_BYTES = 72 // bcrypt max length
)

func checkLength(raw password.RawPassword) (err error) {
	if len(raw) == 0 {
		err = password.ErrCheckLengthEmpty
		return
	}

	if len(raw) > PASSWORD_MAX_BYTES {
		err = password.ErrCheckLengthTooLong
		return
	}

	return nil
}
