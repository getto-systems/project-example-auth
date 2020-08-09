package ticket

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
	infra "github.com/getto-systems/project-example-id/infra/ticket"
)

type (
	action struct {
		logger infra.Logger

		gen credential_infra.TicketNonceGenerator

		tickets infra.TicketRepository
	}
)

func NewAction(
	logger infra.Logger,

	gen credential_infra.TicketNonceGenerator,

	tickets infra.TicketRepository,
) action {
	return action{
		logger: logger,

		gen: gen,

		tickets: tickets,
	}
}
