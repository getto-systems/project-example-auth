package ticket_core

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

// user が正しいことは確認済みでなければならない
func (action action) Validate(request request.Request, user user.User, ticket credential.Ticket) (err error) {
	action.logger.TryToValidate(request, user, ticket.Nonce())

	dataUser, expires, found, err := action.tickets.FindUserAndExpires(ticket.Nonce())
	if err != nil {
		action.logger.FailedToValidate(request, user, ticket.Nonce(), err)
		return
	}
	if !found {
		err = errValidateNotFoundTicket
		action.logger.FailedToValidateBecauseTicketNotFound(request, user, ticket.Nonce(), err)
		return
	}
	if user.ID() != dataUser.ID() {
		err = errValidateMatchFailedUser
		action.logger.FailedToValidateBecauseUserMatchFailed(request, user, ticket.Nonce(), err)
		return
	}
	if request.Expired(expires) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateBecauseExpired(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.AuthByTicket(request, dataUser, ticket.Nonce())
	return nil
}
