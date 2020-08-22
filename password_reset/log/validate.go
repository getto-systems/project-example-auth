package password_reset_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/password_reset/infra"

	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) validate() infra.ValidateLogger {
	return log
}

func (log Logger) TryToValidateToken(request request.Request, login user.Login) {
	log.logger.Debug(validateLog("TryToValidateToken", request, login, nil, nil))
}
func (log Logger) FailedToValidateToken(request request.Request, login user.Login, err error) {
	log.logger.Error(validateLog("FailedToValidateToken", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionNotFound(request request.Request, login user.Login, err error) {
	log.logger.Audit(validateLog("FailedToValidateTokenBecauseSessionNotFound", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionClosed(request request.Request, login user.Login, err error) {
	log.logger.Info(validateLog("FailedToValidateTokenBecauseSessionClosed", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseSessionExpired(request request.Request, login user.Login, err error) {
	log.logger.Info(validateLog("FailedToValidateTokenBecauseSessionExpired", request, login, nil, err))
}
func (log Logger) FailedToValidateTokenBecauseLoginMatchFailed(request request.Request, login user.Login, err error) {
	log.logger.Audit(validateLog("FailedToValidateTokenBecauseLoginMatchFailed", request, login, nil, err))
}
func (log Logger) AuthByToken(request request.Request, login user.Login, user user.User) {
	log.logger.Audit(validateLog("AuthByToken", request, login, &user, nil))
}

type (
	validateEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		Login   user.LoginLog      `json:"login"`
		User    *user.UserLog      `json:"user,omitempty"`
		Err     *string            `json:"error,omitempty"`
	}
)

func validateLog(message string, request request.Request, login user.Login, user *user.User, err error) validateEntry {
	entry := validateEntry{
		Action:  "PasswordReset/Validate",
		Message: message,
		Request: request.Log(),
		Login:   login.Log(),
	}

	if user != nil {
		log := user.Log()
		entry.User = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry validateEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
