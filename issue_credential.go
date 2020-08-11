package _usecase

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
)

func (h Backend) issueCredential(request request.Request, ticket credential.Ticket) (_ credential.Credential, err error) {
	ticketToken, err := h.credential.IssueTicket(request, ticket)
	if err != nil {
		return
	}

	apiToken, err := h.credential.IssueApiToken(request, ticket)
	if err != nil {
		return
	}

	contentToken, err := h.credential.IssueContentToken(request, ticket)
	if err != nil {
		return
	}

	return credential.NewCredential(ticketToken, apiToken, contentToken), nil
}
