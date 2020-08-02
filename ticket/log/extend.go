package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) extend() ticket.ExtendLogger {
	return log
}

func (log Logger) TryToExtendTicket(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires) {
	log.logger.Debug(extendEntry("TryToExtendTicket", request, user, nonce, expires, nil))
}
func (log Logger) FailedToExtendTicket(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires, err error) {
	log.logger.Error(extendEntry("FailedToExtendTicket", request, user, nonce, expires, err))
}
func (log Logger) ExtendTicket(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires) {
	log.logger.Error(extendEntry("ExtendTicket", request, user, nonce, expires, nil))
}

func extendEntry(event string, request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Extend/%s", event),
		Request: request,
		User:    &user,
		Nonce:   &nonce,
		Expires: &expires,
		Error:   err,
	}
}
