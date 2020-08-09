package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) createSession() infra.CreateSessionLogger {
	return log
}

func (log Logger) TryToCreateSession(request request.Request, user user.User, login user.Login, expires expiration.Expires) {
	log.logger.Debug(createSessionEntry("TryToCreateSession", request, user, login, expires, nil, nil, nil))
}
func (log Logger) FailedToCreateSession(request request.Request, user user.User, login user.Login, expires expiration.Expires, err error) {
	log.logger.Error(createSessionEntry("FailedToCreateSession", request, user, login, expires, nil, nil, err))
}
func (log Logger) FailedToCreateSessionBecauseDestinationNotFound(request request.Request, user user.User, login user.Login, expires expiration.Expires, err error) {
	log.logger.Info(createSessionEntry("FailedToCreateSessionBecauseDestinationNotFound", request, user, login, expires, nil, nil, err))
}
func (log Logger) CreateSession(request request.Request, user user.User, login user.Login, expires expiration.Expires, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(createSessionEntry("CreateSession", request, user, login, expires, &session, &dest, nil))
}

func createSessionEntry(event string, request request.Request, user user.User, login user.Login, expires expiration.Expires, session *password_reset.Session, dest *password_reset.Destination, err error) log.Entry {
	return log.Entry{
		Message:             fmt.Sprintf("PasswordReset/CreateSession/%s", event),
		Request:             request,
		User:                &user,
		Login:               &login,
		ResetSessionExpires: &expires,

		ResetSession:     session,
		ResetDestination: dest,

		Error: err,
	}
}
