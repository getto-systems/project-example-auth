package credential_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) parseTicket() infra.ParseTicketLogger {
	return log
}

func (log Logger) TryToParseTicket(request request.Request) {
	log.logger.Debug(parseTicketEntry("TryToParseTicket", request, nil, nil))
}
func (log Logger) FailedToParseTicket(request request.Request, err error) {
	log.logger.Error(parseTicketEntry("FailedToParseTicket", request, nil, err))
}
func (log Logger) FailedToParseTicketBecauseNonceMatchFailed(request request.Request, err error) {
	log.logger.Audit(parseTicketEntry("FailedToParseTicketBecauseNonceMatchFailed", request, nil, err))
}
func (log Logger) ParseTicket(request request.Request, user user.User) {
	log.logger.Info(parseTicketEntry("ParseTicket", request, &user, nil))
}

func parseTicketEntry(event string, request request.Request, user *user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Credential/ParseTicket/%s", event),
		Request: request,
		User:    user,

		Error: err,
	}
}
