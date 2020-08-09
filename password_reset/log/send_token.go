package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) sendToken() infra.SendTokenLogger {
	return log
}

func (log Logger) TryToSendToken(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Debug(sendTokenEntry("TryToSendToken", request, session, dest, nil))
}
func (log Logger) FailedToSendToken(request request.Request, session password_reset.Session, dest password_reset.Destination, err error) {
	log.logger.Error(sendTokenEntry("FailedToSendToken", request, session, dest, err))
}
func (log Logger) SendToken(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(sendTokenEntry("SendToken", request, session, dest, nil))
}

func sendTokenEntry(event string, request request.Request, session password_reset.Session, dest password_reset.Destination, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/SendToken/%s", event),
		Request: request,

		PasswordReset: &log.PasswordResetEntry{
			Session:     &session,
			Destination: &dest,
		},

		Error: err,
	}
}
