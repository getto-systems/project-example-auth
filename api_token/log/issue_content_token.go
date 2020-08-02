package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) issue_content_token() api_token.IssueContentTokenLogger {
	return log
}

func (log Logger) TryToIssueContentToken(request request.Request, user user.User, expires time.Expires) {
	log.logger.Debug(issueContentTokenEntry("TryToIssue", request, user, expires, nil))
}
func (log Logger) FailedToIssueContentToken(request request.Request, user user.User, expires time.Expires, err error) {
	log.logger.Error(issueContentTokenEntry("FailedToIssue", request, user, expires, err))
}
func (log Logger) IssueContentToken(request request.Request, user user.User, expires time.Expires) {
	log.logger.Info(issueContentTokenEntry("Issue", request, user, expires, nil))
}

func issueContentTokenEntry(event string, request request.Request, user user.User, expires time.Expires, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("ApiToken/IssueContentToken/%s", event),
		Request: request,
		User:    &user,
		Expires: &expires,
		Error:   err,
	}
}
