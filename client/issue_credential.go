package client

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (h Backend) issueCredential(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires) (_ credential.Credential, err error) {
	ticket, err := h.credential.IssueTicket(request, user, nonce, expires)
	if err != nil {
		return
	}

	apiToken, err := h.credential.IssueApiToken(request, user, expires)
	if err != nil {
		return
	}

	contentToken, err := h.credential.IssueContentToken(request, user, expires)
	if err != nil {
		return
	}

	return credential.NewCredential(ticket, apiToken, contentToken, expires), nil
}
