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

func (log Logger) TryToValidate(request request.Request, nonce ticket.Nonce) {
	log.logger.Debug(validateEntry("TryToValidate", request, nil, nonce, nil))
}
func (log Logger) FailedToValidate(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Error(validateEntry("FailedToValidate", request, nil, nonce, err))
}
func (log Logger) FailedToValidateBecauseTicketNotFound(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseTicketNotFound", request, nil, nonce, err))
}
func (log Logger) FailedToValidateBecauseMatchFailed(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseMatchFailed", request, nil, nonce, err))
}
func (log Logger) FailedToValidateBecauseExpired(request request.Request, nonce ticket.Nonce, err error) {
	log.logger.Info(validateEntry("FailedToValidateBecauseExpired", request, nil, nonce, err))
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
