package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) ParseTicketSignature(request request.Request, nonce credential.TicketNonce, signature credential.TicketSignature) (_ user.User, err error) {
	action.logger.TryToParseTicketSignature(request)

	user, ticketNonce, err := action.ticketParser.Parse(signature)
	if err != nil {
		action.logger.FailedToParseTicketSignature(request, err)
		return
	}
	if ticketNonce != nonce {
		err = credential.ErrParseTicketSignatureMatchFailedNonce
		action.logger.FailedToParseTicketSignatureBecauseNonceMatchFailed(request, err)
		return
	}

	action.logger.ParseTicketSignature(request, user)
	return user, nil
}
