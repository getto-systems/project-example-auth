package password_core

import (
	"github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/password"
)

type (
	action struct {
		logger infra.Logger

		exp ticket.Expiration

		gen     infra.PasswordGenerator
		matcher infra.PasswordMatcher

		passwords infra.PasswordRepository
	}
)

func NewAction(
	logger infra.Logger,

	exp ticket.Expiration,
	enc infra.PasswordEncrypter,

	passwords infra.PasswordRepository,
) password.Action {
	return action{
		logger: logger,

		exp: exp,

		gen:     enc,
		matcher: enc,

		passwords: passwords,
	}
}
