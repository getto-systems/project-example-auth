package ticket

import (
	"github.com/getto-systems/project-example-id/_misc/errors"

	"github.com/getto-systems/project-example-id/credential"
)

var (
	ErrExtendNotFoundNonce = errors.NewError("Ticket.Extend", "NotFound.Nonce")

	ErrValidateNotFoundTicket  = credential.NewClearCredentialError("Ticket.Validate", "NotFound.Ticket")
	ErrValidateMatchFailedUser = credential.NewClearCredentialError("Ticket.Validate", "MatchFailed.User")
	ErrValidateAlreadyExpired  = credential.NewClearCredentialError("Ticket.Validate", "AlreadyExpired")
)
