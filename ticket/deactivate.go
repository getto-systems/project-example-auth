package ticket

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

type Deactivate struct {
	logger  ticket.DeactivateLogger
	tickets ticket.TicketRepository
}

func NewDeactivate(logger ticket.DeactivateLogger, tickets ticket.TicketRepository) Deactivate {
	return Deactivate{
		logger:  logger,
		tickets: tickets,
	}
}

func (action Deactivate) Deactivate(request request.Request, user user.User, ticket ticket.Ticket) (err error) {
	// user が正しいことは確認済みでなければならない
	action.logger.TryToDeactivate(request, user, ticket.Nonce())

	err = action.tickets.DeactivateExpiresAndExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToDeactivate(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.Deactivate(request, user, ticket.Nonce())
	return nil
}
