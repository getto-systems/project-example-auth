package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) pushSendTokenJob() infra.PushSendTokenJobLogger {
	return log
}

func (log Logger) TryToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Debug(pushSendTokenJobLog("TryToPushSendTokenJob", request, session, dest, nil))
}
func (log Logger) FailedToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination, err error) {
	log.logger.Error(pushSendTokenJobLog("FailedToPushSendTokenJob", request, session, dest, err))
}
func (log Logger) PushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(pushSendTokenJobLog("PushSendTokenJob", request, session, dest, nil))
}

type (
	pushSendTokenJobEntry struct {
		Action      string                        `json:"action"`
		Message     string                        `json:"message"`
		Request     request.RequestLog            `json:"request"`
		Session     password_reset.SessionLog     `json:"session"`
		Destination password_reset.DestinationLog `json:"destination"`
		Err         *string                       `json:"error,omitempty"`
	}
)

func pushSendTokenJobLog(message string, request request.Request, session password_reset.Session, dest password_reset.Destination, err error) pushSendTokenJobEntry {
	entry := pushSendTokenJobEntry{
		Action:      "PasswordReset/PushSendTokenJob",
		Message:     message,
		Request:     request.Log(),
		Session:     session.Log(),
		Destination: dest.Log(),
	}

	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry pushSendTokenJobEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
