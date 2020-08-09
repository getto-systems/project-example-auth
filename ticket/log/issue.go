package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) issue() ticket.IssueLogger {
	return log
}

func (log Logger) TryToIssue(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit) {
	log.logger.Debug(issueEntry("TryToIssue", request, user, expires, limit, nil, nil))
}
func (log Logger) FailedToIssue(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, err error) {
	log.logger.Error(issueEntry("FailedToIssue", request, user, expires, limit, nil, err))
}
func (log Logger) Issue(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, nonce api_token.TicketNonce) {
	log.logger.Info(issueEntry("Issue", request, user, expires, limit, &nonce, nil))
}

func issueEntry(event string, request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, nonce *api_token.TicketNonce, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Ticket/Issue/%s", event),
		Request:     request,
		User:        &user,
		Expires:     &expires,
		ExtendLimit: &limit,
		TicketNonce: nonce,
		Error:       err,
	}
}
