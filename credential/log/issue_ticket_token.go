package credential_log

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) issueTicketToken() infra.IssueTicketTokenLogger {
	return log
}

func (log Logger) TryToIssueTicketToken(request request.Request, user user.User, expires credential.TicketExpires) {
	log.logger.Debug(issueTicketTokenLog("TryToIssue", request, user, expires, nil))
}
func (log Logger) FailedToIssueTicketToken(request request.Request, user user.User, expires credential.TicketExpires, err error) {
	log.logger.Error(issueTicketTokenLog("FailedToIssue", request, user, expires, err))
}
func (log Logger) IssueTicketToken(request request.Request, user user.User, expires credential.TicketExpires) {
	log.logger.Info(issueTicketTokenLog("Issue", request, user, expires, nil))
}

type (
	issueTicketTokenEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Expires string             `json:"expires"`
		Err     *string            `json:"error,omitempty"`
	}
)

func issueTicketTokenLog(message string, request request.Request, user user.User, expires credential.TicketExpires, err error) issueTicketTokenEntry {
	entry := issueTicketTokenEntry{
		Action:  "Credential/IssueTicketToken",
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

func (entry issueTicketTokenEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
