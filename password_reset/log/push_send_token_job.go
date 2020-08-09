package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/gateway/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
)

func (log Logger) pushSendTokenJob() infra.PushSendTokenJobLogger {
	return log
}

func (log Logger) TryToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Debug(pushSendTokenJobEntry("TryToPushSendTokenJob", request, session, dest, nil))
}
func (log Logger) FailedToPushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination, err error) {
	log.logger.Error(pushSendTokenJobEntry("FailedToPushSendTokenJob", request, session, dest, err))
}
func (log Logger) PushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(pushSendTokenJobEntry("PushSendTokenJob", request, session, dest, nil))
}

func pushSendTokenJobEntry(event string, request request.Request, session password_reset.Session, dest password_reset.Destination, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/PushSendTokenJob/%s", event),
		Request: request,

		PasswordReset: &log.PasswordResetEntry{
			Session:     &session,
			Destination: &dest,
		},

		Error: err,
	}
}
