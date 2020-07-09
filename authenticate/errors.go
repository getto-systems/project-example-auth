package authenticate

import (
	"errors"
)

var (
	ErrTicketAuthFailed   = errors.New("ticket-authenticate failed")
	ErrPasswordAuthFailed = errors.New("password-authenticate failed")
)
