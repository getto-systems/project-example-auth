package password_core

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/password"
)

type (
	action struct {
		logger infra.Logger

		extendSecond expiration.ExtendSecond

		generator infra.PasswordGenerator
		matcher   infra.PasswordMatcher

		passwords infra.PasswordRepository
	}
)

func NewAction(
	logger infra.Logger,

	extendSecond expiration.ExtendSecond,
	encrypter infra.PasswordEncrypter,

	passwords infra.PasswordRepository,
) password.Action {
	return action{
		logger: logger,

		extendSecond: extendSecond,

		generator: encrypter,
		matcher:   encrypter,

		passwords: passwords,
	}
}
