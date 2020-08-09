package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) extend() infra.ExtendLogger {
	return log
}

func (log Logger) TryToExtend(request request.Request, user user.User, nonce credential.TicketNonce) {
	log.logger.Debug(extendEntry("TryToExtend", request, user, nonce, nil, nil, nil))
}
func (log Logger) FailedToExtend(request request.Request, user user.User, nonce credential.TicketNonce, err error) {
	log.logger.Error(extendEntry("FailedToExtend", request, user, nonce, nil, nil, err))
}
func (log Logger) Extend(request request.Request, user user.User, nonce credential.TicketNonce, expires time.Expires, limit time.ExtendLimit) {
	log.logger.Info(extendEntry("Extend", request, user, nonce, &expires, &limit, nil))
}

func extendEntry(event string, request request.Request, user user.User, nonce credential.TicketNonce, expires *time.Expires, limit *time.ExtendLimit, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Ticket/Extend/%s", event),
		Request:     request,
		User:        &user,
		TicketNonce: &nonce,
		Expires:     expires,
		ExtendLimit: limit,
		Error:       err,
	}
}
