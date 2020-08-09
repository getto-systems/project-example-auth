package password_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User) {
	log.logger.Debug(validateEntry("TryToValidate", request, user, nil))
}
func (log Logger) FailedToValidateBecausePasswordCheckFailed(request request.Request, user user.User, err error) {
	log.logger.Info(validateEntry("FailedToValidateBecausePasswordCheckFailed", request, user, err))
}
func (log Logger) FailedToValidateBecausePasswordNotFound(request request.Request, user user.User, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecausePasswordNotFound", request, user, err))
}
func (log Logger) FailedToValidateBecausePasswordMatchFailed(request request.Request, user user.User, err error) {
	log.logger.Audit(validateEntry("FailedToValidateBecausePasswordMatchFailed", request, user, err))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, err error) {
	log.logger.Error(validateEntry("FailedToValidate", request, user, err))
}
func (log Logger) AuthByPassword(request request.Request, user user.User) {
	log.logger.Audit(validateEntry("AuthByPassword", request, user, nil))
}

func validateEntry(event string, request request.Request, user user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Password/Validate/%s", event),
		Request: request,
		User:    &user,
		Error:   err,
	}
}
