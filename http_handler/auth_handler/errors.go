package auth_handler

import (
	"errors"
)

var (
	ErrSignedTicketCookieNotFound = errors.New("signed ticket cookie not found")
)
