package ticket_core

import (
	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/ticket"
)

type (
	action struct {
		logger infra.Logger

		gen infra.TicketNonceGenerator

		tickets infra.TicketRepository
	}
)

func NewAction(
	logger infra.Logger,

	gen infra.TicketNonceGenerator,

	tickets infra.TicketRepository,
) ticket.Action {
	return action{
		logger: logger,

		gen: gen,

		tickets: tickets,
	}
}
