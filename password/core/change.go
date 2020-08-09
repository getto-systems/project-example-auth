package password_core

import (
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (action action) Change(request request.Request, user user.User, raw password.RawPassword) (err error) {
	action.logger.TryToChange(request, user)

	err = checkLength(raw)
	if err != nil {
		action.logger.FailedToChangeBecausePasswordCheckFailed(request, user, err)
		return
	}

	hashed, err := action.generator.GeneratePassword(raw)
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
