package api_token

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type IssueApiToken struct {
	logger   api_token.IssueApiTokenLogger
	signer   api_token.ApiTokenSigner
	apiUsers api_token.ApiUserRepository
}

func NewIssueApiToken(logger api_token.IssueApiTokenLogger, signer api_token.ApiTokenSigner, apiUsers api_token.ApiUserRepository) IssueApiToken {
	return IssueApiToken{
		logger:   logger,
		signer:   signer,
		apiUsers: apiUsers,
	}
}

func (action IssueApiToken) Issue(request request.Request, user user.User, expires time.Expires) (_ api_token.ApiToken, err error) {
	action.logger.TryToIssueApiToken(request, user, expires)

	roles, found, err := action.apiUsers.FindApiRoles(user)
	if err != nil {
		action.logger.FailedToIssueApiToken(request, user, expires, err)
		return
	}
	if !found {
		// 見つからない場合は「権限なし」でトークンを発行
		roles = api_token.EmptyApiRoles()
	}

	token, err := action.signer.Sign(user, roles, expires)
	if err != nil {
		action.logger.FailedToIssueApiToken(request, user, expires, err)
		return
	}

	action.logger.IssueApiToken(request, user, roles, expires)
	return token, nil
}
