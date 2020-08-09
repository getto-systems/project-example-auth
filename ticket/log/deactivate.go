package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) deactivate() ticket_infra.DeactivateLogger {
	return log
}

func (log Logger) TryToDeactivate(request request.Request, user user.User, nonce credential.TicketNonce) {
	log.logger.Debug(deactivateEntry("TryToDeactivate", request, user, nonce, nil))
}
func (log Logger) FailedToDeactivate(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Error(deactivateEntry("FailedToDeactivate", request, user, nonce, err))
}
func (log Logger) Deactivate(request request.Request, user user.User, nonce credential.TicketNonce) {
	log.logger.Info(deactivateEntry("Deactivate", request, user, nonce, nil))
}

func deactivateEntry(event string, request request.Request, user user.User, nonce credential.TicketNonce, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Ticket/Deactivate/%s", event),
		Request:     request,
		User:        &user,
		TicketNonce: &nonce,
		Error:       err,
	}
}
