package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) deactivate() infra.DeactivateLogger {
	return log
}

func (log Logger) TryToDeactivate(request request.Request, user user.User) {
	log.logger.Debug(deactivateEntry("TryToDeactivate", request, user, nil))
}
func (log Logger) FailedToDeactivate(request request.Request, user user.User, err error) {
	log.logger.Error(deactivateEntry("FailedToDeactivate", request, user, err))
}
func (log Logger) Deactivate(request request.Request, user user.User) {
	log.logger.Info(deactivateEntry("Deactivate", request, user, nil))
}

func deactivateEntry(event string, request request.Request, user user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Deactivate/%s", event),
		Request: request,
		User:    &user,

		Error: err,
	}
}
