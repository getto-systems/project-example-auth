package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) getUser() user.GetUserLogger {
	return log
}

func (log Logger) TryToGetUser(request request.Request, login user.Login) {
	log.logger.Debug(getUserEntry("TryToGetUser", request, login, nil, nil))
}
func (log Logger) FailedToGetUser(request request.Request, login user.Login, err error) {
	log.logger.Info(getUserEntry("FailedToGetUser", request, login, nil, err))
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
