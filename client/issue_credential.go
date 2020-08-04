package client

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (h Backend) issueCredential(request request.Request, user user.User) (_ data.Credential, err error) {
	newTicket, expires, err := h.ticket.issue.Issue(request, user)
	if err != nil {
		return
	}

	apiToken, contentToken, err := h.issueApiToken(request, user, expires)
	if err != nil {
		return
	}

	return data.NewCredential(newTicket, apiToken, contentToken, expires), nil
}

func (h Backend) issueCredentialByTicket(request request.Request, user user.User, ticket ticket.Ticket, expires time.Expires) (_ data.Credential, err error) {
	apiToken, contentToken, err := h.issueApiToken(request, user, expires)
	if err != nil {
		return
	}

	return data.NewCredential(ticket, apiToken, contentToken, expires), nil
}

func (h Backend) issueApiToken(request request.Request, user user.User, expires time.Expires) (_ api_token.ApiToken, _ api_token.ContentToken, err error) {
	apiToken, err := h.apiToken.issueApiToken.Issue(request, user, expires)
	if err != nil {
		return
	}

	contentToken, err := h.apiToken.issueContentToken.Issue(request, user, expires)
	if err != nil {
		return
	}

	return apiToken, contentToken, nil
}