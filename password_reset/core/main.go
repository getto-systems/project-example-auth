package password_reset_core

import (
	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	action struct {
		logger infra.Logger

		ticketExp ticket.Expiration
		exp       password_reset.Expiration
		gen       infra.SessionGenerator

		sessions     infra.SessionRepository
		destinations infra.DestinationRepository

		tokenQueue infra.SendTokenJobQueue

		tokenSender infra.TokenSender
	}
)

func NewAction(
	logger infra.Logger,

	ticketExp ticket.Expiration,
	exp password_reset.Expiration,
	gen infra.SessionGenerator,

	sessions infra.SessionRepository,
	destinations infra.DestinationRepository,

	tokenQueue infra.SendTokenJobQueue,

	tokenSender infra.TokenSender,
) password_reset.Action {
	return action{
		logger: logger,

		ticketExp: ticketExp,
		exp:       exp,
		gen:       gen,

		sessions:     sessions,
		destinations: destinations,

		tokenQueue: tokenQueue,

		tokenSender: tokenSender,
	}
}
