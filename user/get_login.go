package user

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errGetLoginNotFoundLogin = data.NewError("User.GetLogin", "NotFound.Login")
)

type GetLogin struct {
	logger user.GetLoginLogger
	users  user.UserRepository
}

func NewGetLogin(logger user.GetLoginLogger, users user.UserRepository) GetLogin {
	return GetLogin{
		logger: logger,
		users:  users,
	}
}

func (action GetLogin) Get(request request.Request, user user.User) (_ user.Login, err error) {
	action.logger.TryToGetLogin(request, user)

	login, found, err := action.users.FindLogin(user)
	if err != nil {
		action.logger.FailedToGetLogin(request, user, err)
		return
	}
	if !found {
		err = errGetLoginNotFoundLogin
		action.logger.FailedToGetLogin(request, user, err)
		return
	}

	action.logger.GetLogin(request, user, login)
	return login, nil
}
