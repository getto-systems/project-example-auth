package credential_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) issue_credential() infra.IssueApiTokenLogger {
	return log
}

func (log Logger) TryToIssueApiToken(request request.Request, user user.User, expires expiration.Expires) {
	log.logger.Debug(issueApiTokenEntry("TryToIssue", request, user, expires, nil, nil))
}
func (log Logger) FailedToIssueApiToken(request request.Request, user user.User, expires expiration.Expires, err error) {
	log.logger.Error(issueApiTokenEntry("FailedToIssue", request, user, expires, nil, err))
}
func (log Logger) IssueApiToken(request request.Request, user user.User, roles credential.ApiRoles, expires expiration.Expires) {
	log.logger.Info(issueApiTokenEntry("Issue", request, user, expires, &roles, nil))
}

func issueApiTokenEntry(event string, request request.Request, user user.User, expires expiration.Expires, roles *credential.ApiRoles, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Credential/IssueApiToken/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			Expires:  &expires,
			ApiRoles: roles,
		},

		Error: err,
	}
}
