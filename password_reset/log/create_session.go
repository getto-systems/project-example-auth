package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/password_reset/infra"

	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) createSession() infra.CreateSessionLogger {
	return log
}

func (log Logger) TryToCreateSession(request request.Request, user user.User, login user.Login, expires password_reset.Expires) {
	log.logger.Debug(createSessionLog("TryToCreateSession", request, user, login, expires, nil, nil, nil))
}
func (log Logger) FailedToCreateSession(request request.Request, user user.User, login user.Login, expires password_reset.Expires, err error) {
	log.logger.Error(createSessionLog("FailedToCreateSession", request, user, login, expires, nil, nil, err))
}
func (log Logger) FailedToCreateSessionBecauseDestinationNotFound(request request.Request, user user.User, login user.Login, expires password_reset.Expires, err error) {
	log.logger.Info(createSessionLog("FailedToCreateSessionBecauseDestinationNotFound", request, user, login, expires, nil, nil, err))
}
func (log Logger) CreateSession(request request.Request, user user.User, login user.Login, expires password_reset.Expires, session password_reset.Session, dest password_reset.Destination) {
	log.logger.Info(createSessionLog("CreateSession", request, user, login, expires, &session, &dest, nil))
}

type (
	createSessionEntry struct {
		Action      string                         `json:"action"`
		Message     string                         `json:"message"`
		Request     request.RequestLog             `json:"request"`
		User        user.UserLog                   `json:"user"`
		Login       user.LoginLog                  `json:"login"`
		Expires     string                         `json:"expires"`
		Session     *password_reset.SessionLog     `json:"session,omitempty"`
		Destination *password_reset.DestinationLog `json:"destination,omitempty"`
		Err         *string                        `json:"error,omitempty"`
	}
)

func createSessionLog(message string, request request.Request, user user.User, login user.Login, expires password_reset.Expires, session *password_reset.Session, dest *password_reset.Destination, err error) createSessionEntry {
	entry := createSessionEntry{
		Action:  "PasswordReset/CreateSession",
		Message: message,
		Request: request.Log(),
	}

	if session != nil {
		log := session.Log()
		entry.Session = &log
	}
	if dest != nil {
		log := dest.Log()
		entry.Destination = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry createSessionEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
