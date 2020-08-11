package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) closeSession() infra.CloseSessionLogger {
	return log
}

func (log Logger) TryToCloseSession(request request.Request, session password_reset.Session) {
	log.logger.Debug(closeSessionLog("TryToCloseSession", request, session, nil))
}
func (log Logger) FailedToCloseSession(request request.Request, session password_reset.Session, err error) {
	log.logger.Error(closeSessionLog("FailedToCloseSession", request, session, err))
}
func (log Logger) CloseSession(request request.Request, session password_reset.Session) {
	log.logger.Info(closeSessionLog("CloseSession", request, session, nil))
}

type (
	closeSessionEntry struct {
		Action  string                    `json:"action"`
		Message string                    `json:"message"`
		Request request.RequestLog        `json:"request"`
		Session password_reset.SessionLog `json:"session"`
		Err     *string                   `json:"error,omitempty"`
	}
)

func closeSessionLog(message string, request request.Request, session password_reset.Session, err error) closeSessionEntry {
	entry := closeSessionEntry{
		Action:  "PasswordReset/CloseSession",
		Message: message,
		Request: request.Log(),
		Session: session.Log(),
	}

	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry closeSessionEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
