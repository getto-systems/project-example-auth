package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
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
func (log Logger) Extend(request request.Request, user user.User, nonce credential.TicketNonce, expires expiration.Expires, limit expiration.ExtendLimit) {
	log.logger.Info(extendEntry("Extend", request, user, nonce, &expires, &limit, nil))
}

func extendEntry(event string, request request.Request, user user.User, nonce credential.TicketNonce, expires *expiration.Expires, limit *expiration.ExtendLimit, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Extend/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			TicketNonce: &nonce,
			Expires:     expires,
			ExtendLimit: limit,
		},

		Error: err,
	}
}
