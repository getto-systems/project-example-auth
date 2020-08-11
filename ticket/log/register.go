package ticket_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) register() infra.RegisterLogger {
	return log
}

func (log Logger) TryToRegister(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Debug(registerEntry("TryToRegister", request, user, expires, limit, nil))
}
func (log Logger) FailedToRegister(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit, err error) {
	log.logger.Error(registerEntry("FailedToRegister", request, user, expires, limit, err))
}
func (log Logger) Register(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Info(registerEntry("Register", request, user, expires, limit, nil))
}

func registerEntry(event string, request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Ticket/Register/%s", event),
		Request: request,
		User:    &user,

		Credential: &log.CredentialEntry{
			TicketExpires:     &expires,
			TicketExtendLimit: &limit,
		},

		Error: err,
	}
}
