package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() password_reset.ValidateLogger {
	return log
}

func (log Logger) TryToValidateToken(request request.Request, login user.Login) {
	log.logger.Debug(validateEntry("TryToValidateToken", request, login, nil, nil))
}
func (log Logger) FailedToValidateToken(request request.Request, login user.Login, err error) {
	log.logger.Info(validateEntry("FailedToValidateToken", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseForbidden(request request.Request, login user.Login, err error) {
	log.logger.Audit(validateEntry("FailedToValidateTokenBecauseForbidden", request, login, nil, err))
}
func (log Logger) AuthByToken(request request.Request, login user.Login, user user.User) {
	log.logger.Audit(validateEntry("AuthByToken", request, login, &user, nil))
}

func validateEntry(event string, request request.Request, login user.Login, user *user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("PasswordReset/Validate/%s", event),
		Request: request,
		User:    user,
		Login:   &login,
		Error:   err,
	}
}
