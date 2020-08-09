package password_reset

import (
	infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/ticket"
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
