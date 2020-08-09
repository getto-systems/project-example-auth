package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

func (log Logger) closeSession() password_reset_infra.CloseSessionLogger {
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

		ResetSession: &session,

		Error: err,
	}
}
