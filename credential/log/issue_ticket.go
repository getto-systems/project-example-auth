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

func (log Logger) issueTicket() infra.IssueTicketLogger {
	return log
}

func (log Logger) TryToIssueTicket(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires) {
	log.logger.Debug(issueTicketEntry("TryToIssue", request, user, nonce, expires, nil))
}
func (log Logger) FailedToIssueTicket(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires, err error) {
	log.logger.Error(issueTicketEntry("FailedToIssue", request, user, nonce, expires, err))
}
func (log Logger) IssueTicket(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires) {
	log.logger.Info(issueTicketEntry("Issue", request, user, nonce, expires, nil))
}

func issueTicketEntry(event string, request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires, err error) log.Entry {
	return log.Entry{
		Message:           fmt.Sprintf("Credential/IssueTicket/%s", event),
		Request:           request,
		User:              &user,
		CredentialExpires: &expires,
		TicketNonce:       &nonce,
		Error:             err,
	}
}
