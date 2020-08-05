package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

func (log Logger) sendToken() password_reset.SendTokenLogger {
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

		ResetSession:     &session,
		ResetDestination: &dest,

		Error: err,
	}
}
