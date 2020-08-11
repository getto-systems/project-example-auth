package user_core

import (
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) GetLogin(request request.Request, target user.User) (_ user.Login, err error) {
	action.logger.TryToGetLogin(request, target)

	login, found, err := action.users.FindLogin(target)
	if err != nil {
		action.logger.FailedToGetLogin(request, target, err)
		return
	}
	if !found {
		// user には必ず Login が存在するはずなのでログはエラーログ
		err = user.ErrGetLoginNotFoundLogin
		action.logger.FailedToGetLogin(request, target, err)
		return
	}

	action.logger.GetLogin(request, target, login)
	return login, nil
}
