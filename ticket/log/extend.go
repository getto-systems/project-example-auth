package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) extend() infra.ExtendLogger {
	return log
}

func (log Logger) TryToExtend(request request.Request, user user.User) {
	log.logger.Debug(extendEntry("TryToExtend", request, user, nil, nil, nil))
}
func (log Logger) FailedToExtend(request request.Request, user user.User, err error) {
	log.logger.Error(extendEntry("FailedToExtend", request, user, nil, nil, err))
}
func (log Logger) Extend(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Info(extendEntry("Extend", request, user, &expires, &limit, nil))
}

func extendEntry(event string, request request.Request, user user.User, expires *credential.TicketExpires, limit *credential.TicketExtendLimit, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Extend/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			TicketExpires:     expires,
			TicketExtendLimit: limit,
		},

		Error: err,
	}
}
