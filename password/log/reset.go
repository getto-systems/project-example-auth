package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

func (log Logger) reset() password.ResetLogger {
	return log
}

func (log Logger) TryToCreateResetSession(request data.Request, login password.Login, expires data.Expires) {
	log.logger.Debug(resetEntry("TryToCreateResetSession", request, &login, nil, nil, &expires, nil))
}
func (log Logger) FailedToCreateResetSession(request data.Request, login password.Login, expires data.Expires, err error) {
	log.logger.Info(resetEntry("FailedToCreateResetSession", request, &login, nil, nil, &expires, err))
}

func (log Logger) CreatedResetSession(request data.Request, login password.Login, expires data.Expires, user data.User, session password.ResetSession) {
	log.logger.Audit(resetEntry("CreatedResetSession", request, &login, &session, &user, &expires, nil))
}

func (log Logger) TryToGetResetStatus(request data.Request, session password.ResetSession) {
	log.logger.Debug(resetEntry("TryToGetResetStatus", request, nil, &session, nil, nil, nil))
}
func (log Logger) FailedToGetResetStatus(request data.Request, session password.ResetSession, err error) {
	log.logger.Info(resetEntry("FailedToGetResetStatus", request, nil, &session, nil, nil, err))
}

func (log Logger) TryToValidateResetToken(request data.Request) {
	log.logger.Debug(resetEntry("TryToValidateResetToken", request, nil, nil, nil, nil, nil))
}
func (log Logger) FailedToValidateResetToken(request data.Request, err error) {
	log.logger.Debug(resetEntry("FailedToValidateResetToken", request, nil, nil, nil, nil, err))
}

func (log Logger) AuthedByResetToken(request data.Request, user data.User) {
	log.logger.Audit(resetEntry("AuthedByResetToken", request, nil, nil, &user, nil, nil))
}

func resetEntry(
	event string,
	request data.Request,
	login *password.Login,
	session *password.ResetSession,
	user *data.User,
	expires *data.Expires,
	err error,
) event_log.Entry {
	return event_log.Entry{
		Message:      fmt.Sprintf("Password/Reset/%s", event),
		Request:      request,
		Login:        login,
		ResetSession: session,
		User:         user,
		Expires:      expires,
		Error:        err,
	}
}
