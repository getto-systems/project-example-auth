package user

import (
	infra "github.com/getto-systems/project-example-id/infra/user"

	"github.com/getto-systems/project-example-id/data/user"
)

type (
	action struct {
		logger infra.Logger
		users  infra.UserRepository
	}
)

func NewAction(
	logger infra.Logger,
	users infra.UserRepository,
) user.Action {
	return action{
		logger: logger,
		users:  users,
	}
}
