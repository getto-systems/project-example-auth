package user_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/user/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) getLogin() infra.GetLoginLogger {
	return log
}

func (log Logger) TryToGetLogin(request request.Request, user user.User) {
	log.logger.Debug(getLoginEntry("TryToGetLogin", request, user, nil, nil))
}
func (log Logger) FailedToGetLogin(request request.Request, user user.User, err error) {
	log.logger.Error(getLoginEntry("FailedToGetLogin", request, user, nil, err))
}
func (log Logger) GetLogin(request request.Request, user user.User, login user.Login) {
	log.logger.Info(getLoginEntry("GetLogin", request, user, &login, nil))
}

func getLoginEntry(event string, request request.Request, user user.User, login *user.Login, err error) log.Entry {
	return log.Entry{
		Message: fmt.Sprintf("User/GetLogin/%s", event),
		Request: request,
		User:    &user,
		Login:   login,
		Error:   err,
	}
}
