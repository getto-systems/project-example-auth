package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
)

func (action action) IssueApiToken(request request.Request, ticket credential.Ticket) (_ credential.ApiToken, err error) {
	action.logger.TryToIssueApiToken(request, ticket.User(), ticket.TokenExpires())

	roles, found, err := action.apiUsers.FindApiRoles(ticket.User())
	if err != nil {
		action.logger.FailedToIssueApiToken(request, ticket.User(), ticket.TokenExpires(), err)
		return
	}
	if !found {
		// 見つからない場合は「権限なし」でトークンを発行
		roles = credential.EmptyApiRoles()
	}

	signature, err := action.apiTokenSinger.Sign(ticket.User(), roles, ticket.TokenExpires())
	if err != nil {
		action.logger.FailedToIssueApiToken(request, ticket.User(), ticket.TokenExpires(), err)
		return
	}

	token := ticket.NewApiToken(roles, signature)

	action.logger.IssueApiToken(request, ticket.User(), roles, ticket.TokenExpires())
	return token, nil
}
