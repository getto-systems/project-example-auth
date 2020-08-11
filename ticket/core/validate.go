package ticket_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Validate(request request.Request, user user.User, nonce credential.TicketNonce) (err error) {
	action.logger.TryToValidate(request, user)

	dataUser, expires, found, err := action.tickets.FindUserAndExpires(nonce)
	if err != nil {
		action.logger.FailedToValidate(request, user, err)
		return
	}
	if !found {
		err = ticket.ErrValidateNotFoundTicket
		action.logger.FailedToValidateBecauseTicketNotFound(request, user, err)
		return
	}
	if user.ID() != dataUser.ID() {
		err = ticket.ErrValidateMatchFailedUser
		action.logger.FailedToValidateBecauseUserMatchFailed(request, user, err)
		return
	}
	if expires.Expired(request) {
		err = ticket.ErrValidateAlreadyExpired
		action.logger.FailedToValidateBecauseExpired(request, user, err)
		return
	}

	action.logger.AuthByTicket(request, dataUser)
	return nil
}
