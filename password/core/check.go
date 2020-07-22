package core

import (
	"github.com/getto-systems/project-example-id/password"
)

func checkPassword(raw password.RawPassword) error {
	if len(raw) == 0 {
		return password.ErrPasswordIsEmpty
	}
	if len([]byte(raw)) > password.PASSWORD_BYTES_LIMIT {
		return password.ErrPasswordIsTooLong
	}
	return nil
}
