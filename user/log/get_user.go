package user_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/user/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) getUser() infra.GetUserLogger {
	return log
}

func (log Logger) TryToGetUser(request request.Request, login user.Login) {
	log.logger.Debug(getUserLog("TryToGetUser", request, login, nil, nil))
}
func (log Logger) FailedToGetUser(request request.Request, login user.Login, err error) {
	log.logger.Error(getUserLog("FailedToGetUser", request, login, nil, err))
}
func (log Logger) FailedToGetUserBecauseUserNotFound(request request.Request, login user.Login, err error) {
	log.logger.Info(getUserLog("FailedToGetUserBecauseUserNotFound", request, login, nil, err))
}
func (log Logger) GetUser(request request.Request, login user.Login, user user.User) {
	log.logger.Info(getUserLog("GetUser", request, login, &user, nil))
}

type (
	getUserEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		Login   user.LoginLog      `json:"login"`
		User    *user.UserLog      `json:"user,omitempty"`
		Err     *string            `json:"error,omitempty"`
	}
)

func getUserLog(message string, request request.Request, login user.Login, user *user.User, err error) getUserEntry {
	entry := getUserEntry{
		Action:  "User/GetUser",
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

func (entry getUserEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
