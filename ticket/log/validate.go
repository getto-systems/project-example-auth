package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User, nonce credential.TicketNonce) {
	log.logger.Debug(validateEntry("TryToValidate", request, user, nonce, nil))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Error(validateEntry("FailedToValidate", request, user, nonce, err))
}
func (log Logger) FailedToValidateBecauseTicketNotFound(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseTicketNotFound", request, user, nonce, err))
}
func (log Logger) FailedToValidateBecauseUserMatchFailed(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseUserMatchFailed", request, user, nonce, err))
}
func (log Logger) FailedToValidateBecauseExpired(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Info(validateEntry("FailedToValidateBecauseExpired", request, user, nonce, err))
}
func (log Logger) AuthByTicket(request request.Request, user user.User, nonce credential.TicketNonce) {
	log.logger.Info(validateEntry("AuthByTicket", request, user, nonce, nil))
}

func validateEntry(event string, request request.Request, user user.User, nonce credential.TicketNonce, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Ticket/Validate/%s", event),
		Request:     request,
		User:        &user,
		TicketNonce: &nonce,
		Error:       err,
	}
}
