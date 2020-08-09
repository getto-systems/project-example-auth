package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/errors"
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
	action.logger.TryToParseTicket(request, ticket.Nonce())

	user, nonce, err := action.ticketParser.Parse(ticket.Signature())
	if err != nil {
		action.logger.FailedToParseTicket(request, ticket.Nonce(), err)
		return
	}
	if nonce != ticket.Nonce() {
		err = errValidateMatchFailedNonce
		action.logger.FailedToParseTicketBecauseNonceMatchFailed(request, ticket.Nonce(), err)
		return
	}

	action.logger.ParseTicket(request, ticket.Nonce(), user)
	return user, nil
}
