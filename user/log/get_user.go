package user_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/user/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) getUser() infra.GetUserLogger {
	return log
}

func (log Logger) TryToGetUser(request request.Request, login user.Login) {
	log.logger.Debug(getUserEntry("TryToGetUser", request, login, nil, nil))
}
func (log Logger) FailedToGetUser(request request.Request, login user.Login, err error) {
	log.logger.Error(getUserEntry("FailedToGetUser", request, login, nil, err))
}
func (log Logger) FailedToGetUserBecauseUserNotFound(request request.Request, login user.Login, err error) {
	log.logger.Info(getUserEntry("FailedToGetUserBecauseUserNotFound", request, login, nil, err))
}
func (log Logger) GetUser(request request.Request, login user.Login, user user.User) {
	log.logger.Info(getUserEntry("GetUser", request, login, &user, nil))
}

func getUserEntry(event string, request request.Request, login user.Login, user *user.User, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("User/GetUser/%s", event),
		Request: request,
		User:    user,
		Login:   &login,
		Error:   err,
	}
}
