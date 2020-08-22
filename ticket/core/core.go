package ticket_core

import (
	"github.com/getto-systems/project-example-auth/ticket/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/ticket"
)

type (
	action struct {
		logger infra.Logger

		ticketExpireSecond credential.TicketExpireSecond
		tokenExpireSecond  credential.TokenExpireSecond

		nonceGenerator infra.TicketNonceGenerator

		tickets infra.TicketRepository
	}
)

func NewAction(
	logger infra.Logger,

	ticketExpireSecond credential.TicketExpireSecond,
	tokenExpireSecond credential.TokenExpireSecond,

	nonceGenerator infra.TicketNonceGenerator,

	tickets infra.TicketRepository,
) ticket.Action {
	return action{
		logger: logger,

		ticketExpireSecond: ticketExpireSecond,
		tokenExpireSecond:  tokenExpireSecond,

		nonceGenerator: nonceGenerator,

		tickets: tickets,
	}
}
