package credential_core

import (
	"github.com/getto-systems/project-example-id/misc/errors"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errValidateMatchFailedNonce = errors.NewError("Ticket.Validate", "MatchFailed.Nonce")
	errValidateNotFoundTicket   = errors.NewError("Ticket.Validate", "NotFound.Ticket")
	errValidateMatchFailedUser  = errors.NewError("Ticket.Validate", "MatchFailed.User")
	errValidateAlreadyExpired   = errors.NewError("Ticket.Validate", "AlreadyExpired")
)

func (action action) ParseTicket(request request.Request, ticket credential.Ticket) (_ user.User, err error) {
	action.logger.TryToParseTicket(request)

	user, nonce, err := action.ticketParser.Parse(ticket.Signature())
	if err != nil {
		action.logger.FailedToParseTicket(request, err)
		return
	}
	if nonce != ticket.Nonce() {
		err = errValidateMatchFailedNonce
		action.logger.FailedToParseTicketBecauseNonceMatchFailed(request, err)
		return
	}

	action.logger.ParseTicket(request, user)
	return user, nil
}
