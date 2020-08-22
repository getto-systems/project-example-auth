package password_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/password/infra"

	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User) {
	log.logger.Debug(validateLog("TryToValidate", request, user, nil))
}
func (log Logger) FailedToValidateBecausePasswordCheckFailed(request request.Request, user user.User, err error) {
	log.logger.Info(validateLog("FailedToValidateBecausePasswordCheckFailed", request, user, err))
}
func (log Logger) FailedToValidateBecausePasswordNotFound(request request.Request, user user.User, err error) {
	log.logger.Audit(validateLog("FailedToValidateBecausePasswordNotFound", request, user, err))
}
func (log Logger) FailedToValidateBecausePasswordMatchFailed(request request.Request, user user.User, err error) {
	log.logger.Audit(validateLog("FailedToValidateBecausePasswordMatchFailed", request, user, err))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, err error) {
	log.logger.Error(validateLog("FailedToValidate", request, user, err))
}
func (log Logger) AuthByPassword(request request.Request, user user.User) {
	log.logger.Audit(validateLog("AuthByPassword", request, user, nil))
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

func validateLog(message string, request request.Request, user user.User, err error) changeEntry {
	entry := changeEntry{
		Action:  "Password/Validate",
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
