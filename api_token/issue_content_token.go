package api_token

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type IssueContentToken struct {
	logger api_token.IssueContentTokenLogger
	signer api_token.ContentTokenSigner
}

func NewIssueContentToken(logger api_token.IssueContentTokenLogger, signer api_token.ContentTokenSigner) IssueContentToken {
	return IssueContentToken{
		logger: logger,
		signer: signer,
	}
}

func (action IssueContentToken) Issue(request request.Request, user user.User, expires time.Expires) (_ api_token.ContentToken, err error) {
	action.logger.TryToIssueContentToken(request, user, expires)

	token, err := action.signer.Sign(expires)
	if err != nil {
		action.logger.FailedToIssueContentToken(request, user, expires, err)
		return
	}

	action.logger.IssueContentToken(request, user, expires)
	return token, nil
}
