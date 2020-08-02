package user

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errGetUserNotFoundUser = errors.NewError("Password.GetUser", "NotFound.User")
)

type GetUser struct {
	logger user.GetUserLogger
	users  user.UserRepository
}

func NewGetUser(logger user.GetUserLogger, users user.UserRepository) GetUser {
	return GetUser{
		logger: logger,
		users:  users,
	}
}

func (action GetUser) Get(request request.Request, login user.Login) (_ user.User, err error) {
	action.logger.TryToGetUser(request, login)

	user, found, err := action.users.FindUser(login)
	if err != nil {
		action.logger.FailedToGetUser(request, login, err)
		return
	}
	if !found {
		err = errGetUserNotFoundUser
		action.logger.FailedToGetUser(request, login, err)
		return
	}

	action.logger.GetUser(request, login, user)
	return user, nil
}
