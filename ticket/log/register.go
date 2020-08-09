package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) register() ticket_infra.RegisterLogger {
	return log
}

func (log Logger) TryToRegister(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit) {
	log.logger.Debug(registerEntry("TryToRegister", request, user, expires, limit, nil, nil))
}
func (log Logger) FailedToRegister(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, err error) {
	log.logger.Error(registerEntry("FailedToRegister", request, user, expires, limit, nil, err))
}
func (log Logger) Register(request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, nonce credential.TicketNonce) {
	log.logger.Info(registerEntry("Register", request, user, expires, limit, &nonce, nil))
}

func registerEntry(event string, request request.Request, user user.User, expires time.Expires, limit time.ExtendLimit, nonce *credential.TicketNonce, err error) log.Entry {
	return log.Entry{
		Message:     fmt.Sprintf("Ticket/Register/%s", event),
		Request:     request,
		User:        &user,
		Expires:     &expires,
		ExtendLimit: &limit,
		TicketNonce: nonce,
		Error:       err,
	}
}
