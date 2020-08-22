package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/ticket/infra"

	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) deactivate() infra.DeactivateLogger {
	return log
}

func (log Logger) TryToDeactivate(request request.Request, user user.User) {
	log.logger.Debug(deactivateLog("TryToDeactivate", request, user, nil))
}
func (log Logger) FailedToDeactivate(request request.Request, user user.User, err error) {
	log.logger.Error(deactivateLog("FailedToDeactivate", request, user, err))
}
func (log Logger) Deactivate(request request.Request, user user.User) {
	log.logger.Info(deactivateLog("Deactivate", request, user, nil))
}

type (
	deactivateEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Err     *string            `json:"error,omitempty"`
	}
)

func deactivateLog(message string, request request.Request, user user.User, err error) deactivateEntry {
	entry := deactivateEntry{
		Action:  "Ticket/Deactivate",
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

func (entry deactivateEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
