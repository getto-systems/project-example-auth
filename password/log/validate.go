package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

func (log Logger) validate() password.ValidateLogger {
	return log
}

func (log Logger) TryToValidate(request data.Request, login password.Login) {
	log.logger.Debug(validateEntry("TryToValidate", request, login, nil, nil))
}
func (log Logger) FailedToValidate(request data.Request, login password.Login, err error) {
	log.logger.Audit(validateEntry("FailedToValidate", request, login, nil, err))
}

func (log Logger) AuthedByPassword(request data.Request, login password.Login, user data.User) {
	log.logger.Audit(validateEntry("AuthedByPassword", request, login, &user, nil))
}

func validateEntry(event string, request data.Request, login password.Login, user *data.User, err error) event_log.Entry {
	return event_log.Entry{
		Message: fmt.Sprintf("Password/Validate/%s", event),
		Request: request,
		Login:   &login,
		User:    user,
		Error:   err,
	}
}
