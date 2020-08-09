package ticket_core

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/ticket"
)

type (
	action struct {
		logger infra.Logger

		expireSecond   expiration.ExpireSecond
		nonceGenerator infra.TicketNonceGenerator

		tickets infra.TicketRepository
	}
)

func NewAction(
	logger infra.Logger,

	expireSecond expiration.ExpireSecond,
	nonceGenerator infra.TicketNonceGenerator,

	tickets infra.TicketRepository,
) ticket.Action {
	return action{
		logger: logger,

		expireSecond:   expireSecond,
		nonceGenerator: nonceGenerator,

		tickets: tickets,
	}
}
