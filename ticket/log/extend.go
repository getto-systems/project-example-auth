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

func (log Logger) TryToExtend(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires) {
	log.logger.Debug(extendEntry("TryToExtend", request, user, nonce, expires, nil))
}
func (log Logger) FailedToExtend(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires, err error) {
	log.logger.Error(extendEntry("FailedToExtend", request, user, nonce, expires, err))
}
func (log Logger) FailedToExtendBecauseTicketNotFound(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires, err error) {
	log.logger.Audit(extendEntry("FailedToExtendBecauseTicketNotFound", request, user, nonce, expires, err))
}
func (log Logger) FailedToExtendBecauseUserMatchFailed(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires, err error) {
	log.logger.Audit(extendEntry("FailedToExtendBecauseUserMatchFailed", request, user, nonce, expires, err))
}
func (log Logger) Extend(request request.Request, user user.User, nonce ticket.Nonce, expires time.Expires) {
	log.logger.Info(extendEntry("Extend", request, user, nonce, expires, nil))
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
