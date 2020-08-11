package password_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) change() infra.ChangeLogger {
	return log
}

func (log Logger) TryToChange(request request.Request, user user.User) {
	log.logger.Debug(changeLog("TryToChange", request, user, nil))
}
func (log Logger) FailedToChange(request request.Request, user user.User, err error) {
	log.logger.Error(changeLog("FailedToChange", request, user, err))
}
func (log Logger) FailedToChangeBecausePasswordCheckFailed(request request.Request, user user.User, err error) {
	log.logger.Info(changeLog("FailedToChangeBecausePasswordCheckFailed", request, user, err))
}
func (log Logger) Change(request request.Request, user user.User) {
	log.logger.Audit(changeLog("Change", request, user, nil))
}

type (
	changeEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Err     *string            `json:"error,omitempty"`
	}
)

func changeLog(message string, request request.Request, user user.User, err error) changeEntry {
	entry := changeEntry{
		Action:  "Password/Change",
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

func (entry changeEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
