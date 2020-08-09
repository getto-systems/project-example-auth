package password

import (
	infra "github.com/getto-systems/project-example-id/infra/password"

	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type Change struct {
	logger    infra.ChangeLogger
	gen       infra.PasswordGenerator
	passwords infra.PasswordRepository
}

func NewChange(logger infra.ChangeLogger, gen infra.PasswordGenerator, passwords infra.PasswordRepository) Change {
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
		action.logger.FailedToChangeBecausePasswordCheckFailed(request, user, err)
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
