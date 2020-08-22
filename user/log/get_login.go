package user_log

import (
	"fmt"

	"github.com/getto-systems/project-example-auth/user/infra"

	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) getLogin() infra.GetLoginLogger {
	return log
}

func (log Logger) TryToGetLogin(request request.Request, user user.User) {
	log.logger.Debug(getLoginLog("TryToGetLogin", request, user, nil, nil))
}
func (log Logger) FailedToGetLogin(request request.Request, user user.User, err error) {
	log.logger.Error(getLoginLog("FailedToGetLogin", request, user, nil, err))
}
func (log Logger) GetLogin(request request.Request, user user.User, login user.Login) {
	log.logger.Info(getLoginLog("GetLogin", request, user, &login, nil))
}

type (
	getLoginEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    user.UserLog       `json:"user"`
		Login   *user.LoginLog     `json:"login,omitempty"`
		Err     *string            `json:"error,omitempty"`
	}
)

func getLoginLog(message string, request request.Request, user user.User, login *user.Login, err error) getLoginEntry {
	entry := getLoginEntry{
		Action:  "User/GetLogin",
		Message: message,
		Request: request.Log(),
		User:    user.Log(),
	}

	if login != nil {
		log := login.Log()
		entry.Login = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry getLoginEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
