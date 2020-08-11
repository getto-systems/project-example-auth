package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) ParseTicket(request request.Request, nonce credential.TicketNonce, signature credential.TicketSignature) (_ user.User, err error) {
	action.logger.TryToParseTicket(request)

	user, ticketNonce, err := action.ticketParser.Parse(signature)
	if err != nil {
		action.logger.FailedToParseTicket(request, err)
		return
	}
	if ticketNonce != nonce {
		err = credential.ErrParseTicketMatchFailedNonce
		action.logger.FailedToParseTicketBecauseNonceMatchFailed(request, err)
		return
	}

	action.logger.ParseTicket(request, user)
	return user, nil
}
