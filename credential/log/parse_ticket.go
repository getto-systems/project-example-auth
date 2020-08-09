package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) parseTicket() credential.ParseTicketLogger {
	return log
}

func (log Logger) TryToParseTicket(request request.Request, nonce credential.TicketNonce) {
	log.logger.Debug(parseTicketEntry("TryToParseTicket", request, nonce, nil, nil))
}
func (log Logger) FailedToParseTicket(request request.Request, nonce credential.TicketNonce, err error) {
	log.logger.Error(parseTicketEntry("FailedToParseTicket", request, nonce, nil, err))
}
func (log Logger) FailedToParseTicketBecauseNonceMatchFailed(request request.Request, nonce credential.TicketNonce, err error) {
	log.logger.Audit(parseTicketEntry("FailedToParseTicketBecauseNonceMatchFailed", request, nonce, nil, err))
}
func (log Logger) ParseTicket(request request.Request, nonce credential.TicketNonce, user user.User) {
	log.logger.Info(parseTicketEntry("ParseTicket", request, nonce, &user, nil))
}

func parseTicketEntry(event string, request request.Request, nonce credential.TicketNonce, user *user.User, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Credential/ParseTicket/%s", event),
		Request:     request,
		User:        user,
		TicketNonce: &nonce,
		Error:       err,
	}
}
