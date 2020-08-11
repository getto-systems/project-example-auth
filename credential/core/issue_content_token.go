package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
)

func (action action) IssueContentToken(request request.Request, ticket credential.Ticket) (_ credential.ContentToken, err error) {
	action.logger.TryToIssueContentToken(request, ticket.User(), ticket.TokenExpires())

	keyID, policy, signature, err := action.contentTokenSigner.Sign(ticket.TokenExpires())
	if err != nil {
		action.logger.FailedToIssueContentToken(request, ticket.User(), ticket.TokenExpires(), err)
		return
	}

	token := ticket.NewContentToken(keyID, policy, signature)

	action.logger.IssueContentToken(request, ticket.User(), ticket.TokenExpires())
	return token, nil
}
