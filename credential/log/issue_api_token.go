package credential_log

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) issue_credential() infra.IssueApiTokenLogger {
	return log
}

func (log Logger) TryToIssueApiToken(request request.Request, user user.User, expires credential.TokenExpires) {
	log.logger.Debug(issueApiTokenLog("TryToIssue", request, user, expires, nil, nil))
}
func (log Logger) FailedToIssueApiToken(request request.Request, user user.User, expires credential.TokenExpires, err error) {
	log.logger.Error(issueApiTokenLog("FailedToIssue", request, user, expires, nil, err))
}
func (log Logger) IssueApiToken(request request.Request, user user.User, roles credential.ApiRoles, expires credential.TokenExpires) {
	log.logger.Info(issueApiTokenLog("Issue", request, user, expires, &roles, nil))
}

type (
	issueApiTokenEntry struct {
		Action   string             `json:"action"`
		Message  string             `json:"message"`
		Request  request.RequestLog `json:"request"`
		User     user.UserLog       `json:"user"`
		Expires  string             `json:"expires"`
		ApiRoles *[]string          `json:"api_roles,omitempty"`
		Err      *string            `json:"error,omitempty"`
	}
)

func issueApiTokenLog(message string, request request.Request, user user.User, expires credential.TokenExpires, roles *credential.ApiRoles, err error) issueApiTokenEntry {
	entry := issueApiTokenEntry{
		Action:  "Credential/IssueApiToken",
		Message: message,
		Request: request.Log(),
		User:    user.Log(),
		Expires: time.Time(expires).String(),
	}

	if roles != nil {
		apiRoles := []string(*roles)
		entry.ApiRoles = &apiRoles
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry issueApiTokenEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
