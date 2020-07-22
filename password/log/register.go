package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

func (log Logger) register() password.RegisterLogger {
	return log
}

func (log Logger) TryToGetLogin(request data.Request, user data.User) {
	log.logger.Debug(registerEntry("TryToGetLogin", request, user, nil))
}
func (log Logger) FailedToGetLogin(request data.Request, user data.User, err error) {
	log.logger.Info(registerEntry("FailedToGetLogin", request, user, err))
}

func (log Logger) TryToRegister(request data.Request, user data.User) {
	log.logger.Debug(registerEntry("TryToRegister", request, user, nil))
}
func (log Logger) FailedToRegister(request data.Request, user data.User, err error) {
	log.logger.Info(registerEntry("FailedToRegister", request, user, err))
}

func (log Logger) Registered(request data.Request, user data.User) {
	log.logger.Audit(registerEntry("Registered", request, user, nil))
}

func registerEntry(event string, request data.Request, user data.User, err error) event_log.Entry {
	return event_log.Entry{
		Message: fmt.Sprintf("Password/Register/%s", event),
		Request: request,
		User:    &user,
		Error:   err,
	}
}
