package credential

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type IssueApiToken struct {
	logger   credential.IssueApiTokenLogger
	signer   credential.ApiTokenSigner
	apiUsers credential.ApiUserRepository
}

func NewIssueApiToken(logger credential.IssueApiTokenLogger, signer credential.ApiTokenSigner, apiUsers credential.ApiUserRepository) IssueApiToken {
	return IssueApiToken{
		logger:   logger,
		signer:   signer,
		apiUsers: apiUsers,
	}
}

func (action IssueApiToken) Issue(request request.Request, user user.User, expires time.Expires) (_ credential.ApiToken, err error) {
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

	token, err := action.signer.Sign(user, roles, expires)
	if err != nil {
		action.logger.FailedToIssueApiToken(request, user, expires, err)
		return
	}

	action.logger.IssueApiToken(request, user, roles, expires)
	return token, nil
}
