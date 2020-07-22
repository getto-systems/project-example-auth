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

func (log Logger) TryToIssueReset(request data.Request, login password.Login, expires data.Expires) {
	log.logger.Debug(resetEntry("TryToIssueReset", request, &login, nil, nil, &expires, nil))
}
func (log Logger) FailedToIssueReset(request data.Request, login password.Login, expires data.Expires, err error) {
	log.logger.Info(resetEntry("FailedToIssueReset", request, &login, nil, nil, &expires, err))
}

func (log Logger) IssuedReset(request data.Request, reset password.Reset, user data.User, expires data.Expires) {
	log.logger.Audit(resetEntry("IssuedReset", request, nil, &reset, &user, &expires, nil))
}

func (log Logger) TryToGetResetStatus(request data.Request, reset password.Reset) {
	log.logger.Debug(resetEntry("TryToGetResetStatus", request, nil, &reset, nil, nil, nil))
}
func (log Logger) FailedToGetResetStatus(request data.Request, reset password.Reset, err error) {
	log.logger.Info(resetEntry("FailedToGetResetStatus", request, nil, &reset, nil, nil, err))
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
	reset *password.Reset,
	user *data.User,
	expires *data.Expires,
	err error,
) event_log.Entry {
	return event_log.Entry{
		Message: fmt.Sprintf("Password/Reset/%s", event),
		Request: request,
		Login:   login,
		Reset:   reset,
		User:    user,
		Expires: expires,
		Error:   err,
	}
}
