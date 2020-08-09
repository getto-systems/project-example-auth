package ticket

import (
	infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type Deactivate struct {
	logger  infra.DeactivateLogger
	tickets infra.TicketRepository
}

func NewDeactivate(logger infra.DeactivateLogger, tickets infra.TicketRepository) Deactivate {
	return Deactivate{
		logger:  logger,
		tickets: tickets,
	}
}

// user が正しいことは確認済みでなければならない
func (action Deactivate) Deactivate(request request.Request, user user.User, ticket credential.Ticket) (err error) {
	action.logger.TryToDeactivate(request, user, ticket.Nonce())

	err = action.tickets.DeactivateExpiresAndExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToDeactivate(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.Deactivate(request, user, ticket.Nonce())
	return nil
}
