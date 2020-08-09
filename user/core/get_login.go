package user_core

import (
	"github.com/getto-systems/project-example-id/_misc/errors"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errGetLoginNotFoundLogin = errors.NewError("User.GetLogin", "NotFound.Login")
)

func (action action) GetLogin(request request.Request, user user.User) (_ user.Login, err error) {
	action.logger.TryToGetLogin(request, user)

	login, found, err := action.users.FindLogin(user)
	if err != nil {
		action.logger.FailedToGetLogin(request, user, err)
		return
	}
	if !found {
		// user には必ず Login が存在するはずなのでログはエラーログ
		err = errGetLoginNotFoundLogin
		action.logger.FailedToGetLogin(request, user, err)
		return
	}

	action.logger.GetLogin(request, user, login)
	return login, nil
}
