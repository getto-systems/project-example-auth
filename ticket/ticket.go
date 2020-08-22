package ticket

import (
	"errors"
)

var (
	ErrExtendNotFoundNonce = errors.New("Ticket.Extend/NotFound.Nonce")

	ErrValidateNotFoundTicket  = errors.New("Ticket.Validate/NotFound.Ticket")
	ErrValidateMatchFailedUser = errors.New("Ticket.Validate/MatchFailed.User")
	ErrValidateAlreadyExpired  = errors.New("Ticket.Validate/AlreadyExpired")
)
