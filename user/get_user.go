package user

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errGetUserNotFoundUser = data.NewError("User.GetUser", "NotFound.User")
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
