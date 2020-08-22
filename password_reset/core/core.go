package password_reset_core

import (
	"github.com/getto-systems/project-example-auth/password_reset/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/password_reset"
)

type (
	action struct {
		logger infra.Logger

		ticketExtendSecond  credential.TicketExtendSecond
		sessionExpireSecond password_reset.ExpireSecond

		sessionGenerator infra.SessionGenerator

		sessions     infra.SessionRepository
		destinations infra.DestinationRepository

		tokenQueue infra.SendTokenJobQueue

		tokenSender infra.TokenSender
	}
)

func NewAction(
	logger infra.Logger,

	ticketExtendSecond credential.TicketExtendSecond,
	sessionExpireSecond password_reset.ExpireSecond,

	sessionGenerator infra.SessionGenerator,

	sessions infra.SessionRepository,
	destinations infra.DestinationRepository,

	tokenQueue infra.SendTokenJobQueue,

	tokenSender infra.TokenSender,
) password_reset.Action {
	return action{
		logger: logger,

		ticketExtendSecond:  ticketExtendSecond,
		sessionExpireSecond: sessionExpireSecond,

		sessionGenerator: sessionGenerator,

		sessions:     sessions,
		destinations: destinations,

		tokenQueue: tokenQueue,

		tokenSender: tokenSender,
	}
}
