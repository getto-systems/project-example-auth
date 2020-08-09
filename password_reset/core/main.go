package password_reset_core

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	action struct {
		logger infra.Logger

		credentialExtendSecond expiration.ExtendSecond

		exp              password_reset.Expiration
		sessionGenerator infra.SessionGenerator

		sessions     infra.SessionRepository
		destinations infra.DestinationRepository

		tokenQueue infra.SendTokenJobQueue

		tokenSender infra.TokenSender
	}
)

func NewAction(
	logger infra.Logger,

	credentialExtendSecond expiration.ExtendSecond,

	exp password_reset.Expiration,
	sessionGenerator infra.SessionGenerator,

	sessions infra.SessionRepository,
	destinations infra.DestinationRepository,

	tokenQueue infra.SendTokenJobQueue,

	tokenSender infra.TokenSender,
) password_reset.Action {
	return action{
		logger: logger,

		credentialExtendSecond: credentialExtendSecond,

		exp:              exp,
		sessionGenerator: sessionGenerator,

		sessions:     sessions,
		destinations: destinations,

		tokenQueue: tokenQueue,

		tokenSender: tokenSender,
	}
}
