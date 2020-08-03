package ticket

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errShrinkNotFoundNonce = errors.NewError("Ticket.Shrink", "NotFound.Nonce")
	errShrinkDifferentUser = errors.NewError("Ticket.Shrink", "DifferentUser")
)

type Shrink struct {
	logger  ticket.ShrinkLogger
	tickets ticket.TicketRepository
}

func NewShrink(logger ticket.ShrinkLogger, tickets ticket.TicketRepository) Shrink {
	return Shrink{
		logger:  logger,
		tickets: tickets,
	}
}

func (action Shrink) Shrink(request request.Request, user user.User, ticket ticket.Ticket) (err error) {
	action.logger.TryToShrink(request, user, ticket.Nonce())

	ticketUser, found, err := action.tickets.FindUser(ticket.Nonce())
	if err != nil {
		action.logger.FailedToShrink(request, user, ticket.Nonce(), err)
		return
	}
	if !found {
		err = errShrinkNotFoundNonce
		action.logger.FailedToShrink(request, user, ticket.Nonce(), err)
		return
	}
	if ticketUser.ID() != user.ID() {
		err = errShrinkDifferentUser
		action.logger.FailedToShrink(request, user, ticket.Nonce(), err)
		return
	}

	err = action.tickets.ShrinkExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToShrink(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.Shrink(request, user, ticket.Nonce())
	return nil
}
