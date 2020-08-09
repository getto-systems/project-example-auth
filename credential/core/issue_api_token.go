package credential_core

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) IssueApiToken(request request.Request, user user.User, expires expiration.Expires) (_ credential.ApiToken, err error) {
	action.logger.TryToIssueApiToken(request, user, expires)

	roles, found, err := action.apiUsers.FindApiRoles(user)
	if err != nil {
		action.logger.FailedToIssueApiToken(request, user, expires, err)
		return
	}
	if !found {
		// 見つからない場合は「権限なし」でトークンを発行
		roles = credential.EmptyApiRoles()
	}

	token, err := action.apiTokenSinger.Sign(user, roles, expires)
	if err != nil {
		action.logger.FailedToIssueApiToken(request, user, expires, err)
		return
	}

	action.logger.IssueApiToken(request, user, roles, expires)
	return token, nil
}