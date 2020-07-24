package core

import (
	"github.com/getto-systems/project-example-id/password"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
)

func check(raw password.RawPassword) error {
	if len(raw) == 0 {
		return errPasswordEmpty
	}
	if len([]byte(raw)) > PASSWORD_BYTES_LIMIT {
		return errPasswordTooLong
	}
	return nil
}
