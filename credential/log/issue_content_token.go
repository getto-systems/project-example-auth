package credential_log

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) issue_content_token() infra.IssueContentTokenLogger {
	return log
}

func (log Logger) TryToIssueContentToken(request request.Request, user user.User, expires credential.TokenExpires) {
	log.logger.Debug(issueContentTokenLog("TryToIssue", request, user, expires, nil))
}
func (log Logger) FailedToIssueContentToken(request request.Request, user user.User, expires credential.TokenExpires, err error) {
	log.logger.Error(issueContentTokenLog("FailedToIssue", request, user, expires, err))
}
func (log Logger) IssueContentToken(request request.Request, user user.User, expires credential.TokenExpires) {
	log.logger.Info(issueContentTokenLog("Issue", request, user, expires, nil))
}

type (
	issueContentTokenEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Expires string             `json:"expires"`
		Err     *string            `json:"error,omitempty"`
	}
)

func issueContentTokenLog(message string, request request.Request, user user.User, expires credential.TokenExpires, err error) issueContentTokenEntry {
	entry := issueContentTokenEntry{
		Action:  "Credential/IssueContentToken",
		Message: message,
		Request: request.Log(),
		User:    user.Log(),
		Expires: time.Time(expires).String(),
	}

	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry issueContentTokenEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
