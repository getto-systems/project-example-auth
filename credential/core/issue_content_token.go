package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) IssueContentToken(request request.Request, user user.User, expires credential.Expires) (_ credential.ContentToken, err error) {
	action.logger.TryToIssueContentToken(request, user, expires)

	token, err := action.contentTokenSigner.Sign(expires)
	if err != nil {
		action.logger.FailedToIssueContentToken(request, user, expires, err)
		return
	}

	action.logger.IssueContentToken(request, user, expires)
	return token, nil
}
