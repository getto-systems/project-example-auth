package user_core

import (
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) GetUser(request request.Request, login user.Login) (_ user.User, err error) {
	action.logger.TryToGetUser(request, login)

	target, found, err := action.users.FindUser(login)
	if err != nil {
		action.logger.FailedToGetUser(request, login, err)
		return
	}
	if !found {
		err = user.ErrGetUserNotFoundUser
		action.logger.FailedToGetUserBecauseUserNotFound(request, login, err)
		return
	}

	action.logger.GetUser(request, login, target)
	return target, nil
}
