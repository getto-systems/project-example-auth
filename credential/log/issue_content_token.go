package credential_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) issue_content_token() infra.IssueContentTokenLogger {
	return log
}

func (log Logger) TryToIssueContentToken(request request.Request, user user.User, expires expiration.Expires) {
	log.logger.Debug(issueContentTokenEntry("TryToIssue", request, user, expires, nil))
}
func (log Logger) FailedToIssueContentToken(request request.Request, user user.User, expires expiration.Expires, err error) {
	log.logger.Error(issueContentTokenEntry("FailedToIssue", request, user, expires, err))
}
func (log Logger) IssueContentToken(request request.Request, user user.User, expires expiration.Expires) {
	log.logger.Info(issueContentTokenEntry("Issue", request, user, expires, nil))
}

func issueContentTokenEntry(event string, request request.Request, user user.User, expires expiration.Expires, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Credential/IssueContentToken/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			Expires: &expires,
		},

		Error: err,
	}
}
