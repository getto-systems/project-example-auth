package auth_handler

import (
	"errors"
)

var (
	ErrEmptyBody                  = errors.New("empty body")
	ErrBodyParseFailed            = errors.New("body parse failed")
	ErrSignedTicketCookieNotFound = errors.New("signed ticket cookie not found")
)
