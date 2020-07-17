package password

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
)

func checkPassword(password data.RawPassword) error {
	if len(password) == 0 {
		return errors.New("password is empty")
	}
	if len([]byte(password)) > PASSWORD_BYTES_LIMIT {
		return errors.New("password is too long")
	}
	return nil
}
