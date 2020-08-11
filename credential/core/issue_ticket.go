package credential_core

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) IssueTicket(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires) (_ credential.TicketToken, err error) {
	action.logger.TryToIssueTicket(request, user, expires)

	signature, err := action.ticketSigner.Sign(user, nonce, expires)
	if err != nil {
		action.logger.FailedToIssueTicket(request, user, expires, err)
		return
	}

	action.logger.IssueTicket(request, user, expires)
	return credential.NewTicketToken(signature, nonce), nil
}
