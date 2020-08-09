package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User) {
	log.logger.Debug(validateEntry("TryToValidate", request, user, nil))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, err error) {
	log.logger.Error(validateEntry("FailedToValidate", request, user, err))
}
func (log Logger) FailedToValidateBecauseTicketNotFound(request request.Request, user user.User, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseTicketNotFound", request, user, err))
}
func (log Logger) FailedToValidateBecauseUserMatchFailed(request request.Request, user user.User, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecauseUserMatchFailed", request, user, err))
}
func (log Logger) FailedToValidateBecauseExpired(request request.Request, user user.User, err error) {
	log.logger.Info(validateEntry("FailedToValidateBecauseExpired", request, user, err))
}
func (log Logger) AuthByTicket(request request.Request, user user.User) {
	log.logger.Info(validateEntry("AuthByTicket", request, user, nil))
}

func validateEntry(event string, request request.Request, user user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Validate/%s", event),
		Request: request,
		User:    &user,

		Error: err,
	}
}
