package credential_core

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) IssueTicket(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires) (_ credential.Ticket, err error) {
	action.logger.TryToIssueTicket(request, user, nonce, expires)

	signature, err := action.ticketSigner.Sign(user, nonce, expires)
	if err != nil {
		action.logger.FailedToIssueTicket(request, user, nonce, expires, err)
		return
	}

	action.logger.IssueTicket(request, user, nonce, expires)
	return credential.NewTicket(signature, nonce), nil
}
