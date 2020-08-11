package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) sendToken() infra.SendTokenLogger {
	return log
}

func (log Logger) TryToSendToken(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Debug(sendTokenLog("TryToSendToken", request, session, dest, nil))
}
func (log Logger) FailedToSendToken(request request.Request, session password_reset.Session, dest password_reset.Destination, err error) {
	log.logger.Error(sendTokenLog("FailedToSendToken", request, session, dest, err))
}
func (log Logger) SendToken(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(sendTokenLog("SendToken", request, session, dest, nil))
}

type (
	sendTokenEntry struct {
		Action      string                        `json:"action"`
		Message     string                        `json:"message"`
		Request     request.RequestLog            `json:"request"`
		Session     password_reset.SessionLog     `json:"session"`
		Destination password_reset.DestinationLog `json:"destination"`
		Err         *string                       `json:"error,omitempty"`
	}
)

func sendTokenLog(message string, request request.Request, session password_reset.Session, dest password_reset.Destination, err error) sendTokenEntry {
	entry := sendTokenEntry{
		Action:      "PasswordReset/SendToken",
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

func (entry sendTokenEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
