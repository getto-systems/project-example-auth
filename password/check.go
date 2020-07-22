package password

import (
	"errors"
)

const (
	PASSWORD_BYTES_LIMIT = 72 // limit of bcrypt
)

var (
	ErrPasswordIsEmpty   = errors.New("password is empty")
	ErrPasswordIsTooLong = errors.New("password is too long")
)
