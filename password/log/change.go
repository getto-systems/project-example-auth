package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) register() password.ChangeLogger {
	return log
}

func (log Logger) TryToChange(request request.Request, user user.User) {
	log.logger.Debug(registerEntry("TryToChange", request, user, nil))
}
func (log Logger) FailedToChange(request request.Request, user user.User, err error) {
	log.logger.Info(registerEntry("FailedToChange", request, user, err))
}
func (log Logger) FailedToChangeBecausePasswordCheckFailed(request request.Request, user user.User, err error) {
	log.logger.Info(registerEntry("FailedToChangeBecausePasswordCheckFailed", request, user, err))
}
func (log Logger) Change(request request.Request, user user.User) {
	log.logger.Audit(registerEntry("Change", request, user, nil))
}

func registerEntry(event string, request request.Request, user user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Password/Change/%s", event),
		Request: request,
		User:    &user,
		Error:   err,
	}
}
