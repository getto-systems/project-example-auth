package password_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) register() infra.ChangeLogger {
	return log
}

func (log Logger) TryToChange(request request.Request, user user.User) {
	log.logger.Debug(registerEntry("TryToChange", request, user, nil))
}
func (log Logger) FailedToChange(request request.Request, user user.User, err error) {
	log.logger.Error(registerEntry("FailedToChange", request, user, err))
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
