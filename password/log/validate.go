package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() password.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request request.Request, user user.User) {
	log.logger.Debug(validateEntry("TryToValidate", request, user, nil))
}
func (log Logger) FailedToValidate(request request.Request, user user.User, err error) {
	log.logger.Audit(validateEntry("FailedToValidate", request, user, err))
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