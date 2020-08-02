package ticket

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errExtendNotFoundNonce = errors.NewError("Ticket.Extend", "NotFound.Nonce")
	errExtendDifferentUser = errors.NewError("Ticket.Extend", "DifferentUser")
)

type Extend struct {
	logger  ticket.ExtendLogger
	signer  ticket.TicketSigner
	exp     expiration
	tickets ticket.TicketRepository
}

func NewExtend(logger ticket.ExtendLogger, signer ticket.TicketSigner, exp ticket.ExpirationParam, tickets ticket.TicketRepository) Extend {
	return Extend{
		logger:  logger,
		signer:  signer,
		exp:     newExpiration(exp),
		tickets: tickets,
	}
}

func (action Extend) Extend(request request.Request, user user.User, oldTicket ticket.Ticket) (_ ticket.Ticket, _ time.Expires, err error) {
	expires := action.exp.Expires(request)
	action.logger.TryToExtendTicket(request, user, oldTicket.Nonce(), expires)

	ticketUser, limit, found, err := action.tickets.FindUserAndExtendLimit(oldTicket.Nonce())
	if err != nil {
		action.logger.FailedToExtendTicket(request, user, oldTicket.Nonce(), expires, err)
		return
	}
	if !found {
		err = errExtendNotFoundNonce
		action.logger.FailedToExtendTicket(request, user, oldTicket.Nonce(), expires, err)
		return
	}
	if ticketUser.ID() != user.ID() {
		err = errExtendDifferentUser
		action.logger.FailedToExtendTicket(request, user, oldTicket.Nonce(), expires, err)
		return
	}

	expires = expires.Limit(limit)

	token, err := action.signer.Sign(user, oldTicket.Nonce(), expires)
	if err != nil {
		action.logger.FailedToExtendTicket(request, user, oldTicket.Nonce(), expires, err)
		return
	}

	newTicket := ticket.NewTicket(token, oldTicket.Nonce())

	action.logger.ExtendTicket(request, user, newTicket.Nonce(), expires)
	return newTicket, expires, nil
}
