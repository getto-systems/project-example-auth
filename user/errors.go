package user

import (
	"errors"
)

var (
	ErrTicketAlreadyExpired = errors.New("ticket already expired")
	ErrUserPasswordNotFound = errors.New("user password not found")
)
