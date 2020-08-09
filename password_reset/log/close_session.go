package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) closeSession() infra.CloseSessionLogger {
	return log
}

func (log Logger) TryToCloseSession(request request.Request, session password_reset.Session) {
	log.logger.Debug(closeSessionEntry("TryToCloseSession", request, session, nil))
}
func (log Logger) FailedToCloseSession(request request.Request, session password_reset.Session, err error) {
	log.logger.Error(closeSessionEntry("FailedToCloseSession", request, session, err))
}
func (log Logger) CloseSession(request request.Request, session password_reset.Session) {
	log.logger.Info(closeSessionEntry("CloseSession", request, session, nil))
}

func closeSessionEntry(event string, request request.Request, session password_reset.Session, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/CloseSession/%s", event),
		Request: request,

		PasswordReset: &log.PasswordResetEntry{
			Session: &session,
		},

		Error: err,
	}
}
