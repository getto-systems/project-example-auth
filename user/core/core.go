package user_core

import (
	"github.com/getto-systems/project-example-auth/user/infra"

	"github.com/getto-systems/project-example-auth/user"
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
