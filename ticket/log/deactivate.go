package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) deactivate() ticket.DeactivateLogger {
	return log
}

func (log Logger) TryToDeactivate(request request.Request, user user.User, nonce ticket.Nonce) {
	log.logger.Debug(deactivateEntry("TryToDeactivate", request, user, nonce, nil))
}
func (log Logger) FailedToDeactivate(request request.Request, user user.User, nonce ticket.Nonce, err error) {
	log.logger.Error(deactivateEntry("FailedToDeactivate", request, user, nonce, err))
}
func (log Logger) Deactivate(request request.Request, user user.User, nonce ticket.Nonce) {
	log.logger.Info(deactivateEntry("Deactivate", request, user, nonce, nil))
}

func deactivateEntry(event string, request request.Request, user user.User, nonce ticket.Nonce, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Deactivate/%s", event),
		Request: request,
		User:    &user,
		Nonce:   &nonce,
		Error:   err,
	}
}