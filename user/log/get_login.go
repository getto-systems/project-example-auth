package log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

func (log Logger) getLogin() user.GetLoginLogger {
	return log
}

func (log Logger) TryToGetLogin(request request.Request, user user.User) {
	log.logger.Debug(getLoginEntry("TryToGetLogin", request, user, nil, nil))
}
func (log Logger) FailedToGetLogin(request request.Request, user user.User, err error) {
	log.logger.Info(getLoginEntry("FailedToGetLogin", request, user, nil, err))
}
func (log Logger) GetLogin(request request.Request, user user.User, login user.Login) {
	log.logger.Debug(getLoginEntry("GetLogin", request, user, &login, nil))
}

func getLoginEntry(event string, request request.Request, user user.User, login *user.Login, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("Password/GetLogin/%s", event),
		Request: request,
		User:    &user,
		Login:   login,
		Error:   err,
	}
}
