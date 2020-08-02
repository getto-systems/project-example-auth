package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

func (log Logger) pushSendTokenJob() password_reset.PushSendTokenJobLogger {
	return log
}

func (log Logger) TryToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Debug(pushSendTokenJobEntry("TryToPushSendTokenJob", request, session, dest, nil))
}
func (log Logger) FailedToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination, err error) {
	log.logger.Info(pushSendTokenJobEntry("FailedToPushSendTokenJob", request, session, dest, err))
}
func (log Logger) PushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(pushSendTokenJobEntry("PushSendTokenJob", request, session, dest, nil))
}

func pushSendTokenJobEntry(event string, request request.Request, session password_reset.Session, dest password_reset.Destination, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/PushSendTokenJob/%s", event),
		Request: request,

		ResetSession:     &session,
		ResetDestination: &dest,

		Error: err,
	}
}
