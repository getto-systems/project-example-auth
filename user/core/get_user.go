package user_core

import (
	"github.com/getto-systems/project-example-id/_misc/errors"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errGetUserNotFoundUser = errors.NewError("User.GetUser", "NotFound.User")
)

func (action action) GetUser(request request.Request, login user.Login) (_ user.User, err error) {
	action.logger.TryToGetUser(request, login)

	user, found, err := action.users.FindUser(login)
	if err != nil {
		action.logger.FailedToGetUser(request, login, err)
		return
	}
	if !found {
		err = errGetUserNotFoundUser
		action.logger.FailedToGetUserBecauseUserNotFound(request, login, err)
		return
	}

	action.logger.GetUser(request, login, user)
	return user, nil
}
