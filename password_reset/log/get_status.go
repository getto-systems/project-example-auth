package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) getStatus() infra.GetStatusLogger {
	return log
}

func (log Logger) TryToGetStatus(request request.Request, session password_reset.Session) {
	log.logger.Debug(getStatusEntry("TryToGetStatus", request, session, nil, nil))
}
func (log Logger) FailedToGetStatus(request request.Request, session password_reset.Session, err error) {
	log.logger.Error(getStatusEntry("FailedToGetStatus", request, session, nil, err))
}
func (log Logger) FailedToGetStatusBecauseSessionNotFound(request request.Request, session password_reset.Session, err error) {
	log.logger.Audit(getStatusEntry("FailedToGetStatusBecauseSessionNotFound", request, session, nil, err))
}
func (log Logger) FailedToGetStatusBecauseLoginMatchFailed(request request.Request, session password_reset.Session, err error) {
	log.logger.Audit(getStatusEntry("FailedToGetStatusBecauseLoginMatchFailed", request, session, nil, err))
}
func (log Logger) GetStatus(request request.Request, session password_reset.Session, status password_reset.Status) {
	log.logger.Info(getStatusEntry("GetStatus", request, session, &status, nil))
}

func getStatusEntry(event string, request request.Request, session password_reset.Session, status *password_reset.Status, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/GetStatus/%s", event),
		Request: request,

		PasswordReset: &log.PasswordResetEntry{
			Session: &session,
			Status:  status,
		},

		Error: err,
	}
}
