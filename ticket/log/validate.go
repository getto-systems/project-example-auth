package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() ticket.ValidateLogger {
	return log
}

func (log Logger) TryToValidateTicket(request request.Request, nonce ticket.Nonce) {
	log.logger.Debug(validateEntry("TryToValidateTicket", request, nil, nonce, nil))
}
func (log Logger) FailedToValidateTicket(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Error(validateEntry("FailedToValidateTicket", request, nil, nonce, err))
}
func (log Logger) FailedToValidateTicketBecauseExpired(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Info(validateEntry("FailedToValidateTicketBecauseExpired", request, nil, nonce, err))
}
func (log Logger) AuthByTicket(request request.Request, user user.User, nonce ticket.Nonce) {
	log.logger.Info(validateEntry("AuthByTicket", request, &user, nonce, nil))
}

func validateEntry(event string, request request.Request, user *user.User, nonce ticket.Nonce, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Validate/%s", event),
		Request: request,
		User:    user,
		Nonce:   &nonce,
		Error:   err,
	}
}
