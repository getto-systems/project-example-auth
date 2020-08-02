package password

import (
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type Change struct {
	logger    password.ChangeLogger
	gen       password.PasswordGenerator
	passwords password.PasswordRepository
}

func NewChange(logger password.ChangeLogger, gen password.PasswordGenerator, passwords password.PasswordRepository) Change {
	return Change{
		logger:    logger,
		gen:       gen,
		passwords: passwords,
	}
}

func (action Change) Change(request request.Request, user user.User, raw password.RawPassword) (err error) {
	action.logger.TryToChange(request, user)

	err = checkLength(raw)
	if err != nil {
		action.logger.FailedToChange(request, user, err)
		return
	}

	hashed, err := action.gen.GeneratePassword(raw)
	if err != nil {
		action.logger.FailedToChange(request, user, err)
		return
	}

	err = action.passwords.ChangePassword(user, hashed)
	if err != nil {
		action.logger.FailedToChange(request, user, err)
		return
	}

	action.logger.Change(request, user)
	return nil
}
