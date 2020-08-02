package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) shrink() ticket.ShrinkLogger {
	return log
}

func (log Logger) TryToShrinkTicket(request request.Request, user user.User, nonce ticket.Nonce) {
	log.logger.Debug(shrinkEntry("TryToShrinkTicket", request, user, nonce, nil))
}
func (log Logger) FailedToShrinkTicket(request request.Request, user user.User, nonce ticket.Nonce, err error) {
	log.logger.Error(shrinkEntry("FailedToShrinkTicket", request, user, nonce, err))
}
func (log Logger) ShrinkTicket(request request.Request, user user.User, nonce ticket.Nonce) {
	log.logger.Info(shrinkEntry("ShrinkTicket", request, user, nonce, nil))
}

func shrinkEntry(event string, request request.Request, user user.User, nonce ticket.Nonce, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Shrink/%s", event),
		Request: request,
		User:    &user,
		Nonce:   &nonce,
		Error:   err,
	}
}
