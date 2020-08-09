package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidateToken(request request.Request, login user.Login) {
	log.logger.Debug(validateEntry("TryToValidateToken", request, login, nil, nil))
}
func (log Logger) FailedToValidateToken(request request.Request, login user.Login, err error) {
	log.logger.Error(validateEntry("FailedToValidateToken", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionNotFound(request request.Request, login user.Login, err error) {
	log.logger.Audit(validateEntry("FailedToValidateTokenBecauseSessionNotFound", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionClosed(request request.Request, login user.Login, err error) {
	log.logger.Info(validateEntry("FailedToValidateTokenBecauseSessionClosed", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionExpired(request request.Request, login user.Login, err error) {
	log.logger.Info(validateEntry("FailedToValidateTokenBecauseSessionExpired", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseLoginMatchFailed(request request.Request, login user.Login, err error) {
	log.logger.Audit(validateEntry("FailedToValidateTokenBecauseLoginMatchFailed", request, login, nil, err))
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
