package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/password_reset/infra"

	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
)

func (log Logger) getStatus() infra.GetStatusLogger {
	return log
}

func (log Logger) TryToGetStatus(request request.Request, session password_reset.Session) {
	log.logger.Debug(getStatusLog("TryToGetStatus", request, session, nil, nil))
}
func (log Logger) FailedToGetStatus(request request.Request, session password_reset.Session, err error) {
	log.logger.Error(getStatusLog("FailedToGetStatus", request, session, nil, err))
}
func (log Logger) FailedToGetStatusBecauseSessionNotFound(request request.Request, session password_reset.Session, err error) {
	log.logger.Audit(getStatusLog("FailedToGetStatusBecauseSessionNotFound", request, session, nil, err))
}
func (log Logger) FailedToGetStatusBecauseLoginMatchFailed(request request.Request, session password_reset.Session, err error) {
	log.logger.Audit(getStatusLog("FailedToGetStatusBecauseLoginMatchFailed", request, session, nil, err))
}
func (log Logger) GetStatus(request request.Request, session password_reset.Session, status password_reset.Status) {
	log.logger.Info(getStatusLog("GetStatus", request, session, &status, nil))
}

type (
	getStatusEntry struct {
		Action  string                    `json:"action"`
		Message string                    `json:"message"`
		Request request.RequestLog        `json:"request"`
		Session password_reset.SessionLog `json:"session"`
		Status  *password_reset.StatusLog `json:"status,omitempty"`
		Err     *string                   `json:"error,omitempty"`
	}
)

func getStatusLog(message string, request request.Request, session password_reset.Session, status *password_reset.Status, err error) getStatusEntry {
	entry := getStatusEntry{
		Action:  "PasswordReset/GetStatus",
		Message: message,
		Request: request.Log(),
		Session: session.Log(),
	}

	if status != nil {
		log := status.Log()
		entry.Status = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry getStatusEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
