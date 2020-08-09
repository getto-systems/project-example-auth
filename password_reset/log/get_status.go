package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

func (log Logger) getStatus() password_reset_infra.GetStatusLogger {
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

		ResetSession: &session,
		ResetStatus:  status,

		Error: err,
	}
}
