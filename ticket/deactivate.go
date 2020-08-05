package ticket

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errDeactivateNotFoundNonce   = data.NewError("Ticket.Deactivate", "NotFound.Nonce")
	errDeactivateMatchFailedUser = data.NewError("Ticket.Deactivate", "MatchFailed.User")
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
	action.logger.TryToDeactivate(request, user, ticket.Nonce())

	ticketUser, found, err := action.tickets.FindUser(ticket.Nonce())
	if err != nil {
		action.logger.FailedToDeactivate(request, user, ticket.Nonce(), err)
		return
	}
	if !found {
		err = errDeactivateNotFoundNonce
		action.logger.FailedToDeactivate(request, user, ticket.Nonce(), err)
		return
	}
	if ticketUser.ID() != user.ID() {
		err = errDeactivateMatchFailedUser
		action.logger.FailedToDeactivateBecauseUserMatchFailed(request, user, ticket.Nonce(), err)
		return
	}

	err = action.tickets.DeactivateExpiresAndExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToDeactivate(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.Deactivate(request, user, ticket.Nonce())
	return nil
}
