package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) issue_api_token() api_token.IssueApiTokenLogger {
	return log
}

func (log Logger) TryToIssueApiToken(request request.Request, user user.User, expires time.Expires) {
	log.logger.Debug(issueApiTokenEntry("TryToIssue", request, user, expires, nil, nil))
}
func (log Logger) FailedToIssueApiToken(request request.Request, user user.User, expires time.Expires, err error) {
	log.logger.Error(issueApiTokenEntry("FailedToIssue", request, user, expires, nil, err))
}
func (log Logger) IssueApiToken(request request.Request, user user.User, roles api_token.ApiRoles, expires time.Expires) {
	log.logger.Info(issueApiTokenEntry("Issue", request, user, expires, &roles, nil))
}

func issueApiTokenEntry(event string, request request.Request, user user.User, expires time.Expires, roles *api_token.ApiRoles, err error) log.Entry {
	return log.Entry{
		Message:  fmt.Sprintf("ApiToken/IssueApiToken/%s", event),
		Request:  request,
		User:     &user,
		Expires:  &expires,
		ApiRoles: roles,
		Error:    err,
	}
}
