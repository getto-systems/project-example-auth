package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User) {
	log.logger.Debug(validateLog("TryToValidate", request, user, nil))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, err error) {
	log.logger.Error(validateLog("FailedToValidate", request, user, err))
}
func (log Logger) FailedToValidateBecauseTicketNotFound(request request.Request, user user.User, err error) {
	log.logger.Audit(validateLog("FailedToValidateBecauseTicketNotFound", request, user, err))
}
func (log Logger) FailedToValidateBecauseUserMatchFailed(request request.Request, user user.User, err error) {
	log.logger.Audit(validateLog("FailedToValidateBecauseUserMatchFailed", request, user, err))
}
func (log Logger) FailedToValidateBecauseExpired(request request.Request, user user.User, err error) {
	log.logger.Info(validateLog("FailedToValidateBecauseExpired", request, user, err))
}
func (log Logger) AuthByTicket(request request.Request, user user.User) {
	log.logger.Info(validateLog("AuthByTicket", request, user, nil))
}

type (
	validateEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Err     *string            `json:"error,omitempty"`
	}
)

func validateLog(message string, request request.Request, user user.User, err error) validateEntry {
	entry := validateEntry{
		Action:  "Ticket/Validate",
		Message: message,
		Request: request.Log(),
		User:    user.Log(),
	}

	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry validateEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
