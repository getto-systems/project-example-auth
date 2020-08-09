package credential

import (
	infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type IssueContentToken struct {
	logger infra.IssueContentTokenLogger
	signer infra.ContentTokenSigner
}

func NewIssueContentToken(logger infra.IssueContentTokenLogger, signer infra.ContentTokenSigner) IssueContentToken {
	return IssueContentToken{
		logger: logger,
		signer: signer,
	}
}

func (action IssueContentToken) Issue(request request.Request, user user.User, expires time.Expires) (_ credential.ContentToken, err error) {
	action.logger.TryToIssueContentToken(request, user, expires)

	token, err := action.signer.Sign(expires)
	if err != nil {
		action.logger.FailedToIssueContentToken(request, user, expires, err)
		return
	}

	action.logger.IssueContentToken(request, user, expires)
	return token, nil
}
