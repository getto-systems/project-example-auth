package credential_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) issueTicket() infra.IssueTicketLogger {
	return log
}

func (log Logger) TryToIssueTicket(request request.Request, user user.User, expires credential.TicketExpires) {
	log.logger.Debug(issueTicketEntry("TryToIssue", request, user, expires, nil))
}
func (log Logger) FailedToIssueTicket(request request.Request, user user.User, expires credential.TicketExpires, err error) {
	log.logger.Error(issueTicketEntry("FailedToIssue", request, user, expires, err))
}
func (log Logger) IssueTicket(request request.Request, user user.User, expires credential.TicketExpires) {
	log.logger.Info(issueTicketEntry("Issue", request, user, expires, nil))
}

func issueTicketEntry(event string, request request.Request, user user.User, expires credential.TicketExpires, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Credential/IssueTicket/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			TicketExpires: &expires,
		},

		Error: err,
	}
}
