package ticket

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateDifferentNonce = data.NewError("Ticket.Validate", "DifferentNonce")
	errValidateNotFoundTicket = data.NewError("Ticket.Validate", "NotFound.Ticket")
	errValidateDifferentUser  = data.NewError("Ticket.Validate", "DifferentUser")
	errValidateAlreadyExpired = data.NewError("Ticket.Validate", "AlreadyExpired")
)

type Validate struct {
	logger  ticket.ValidateLogger
	parser  ticket.TicketParser
	tickets ticket.TicketRepository
}

func NewValidate(logger ticket.ValidateLogger, parser ticket.TicketParser, tickets ticket.TicketRepository) Validate {
	return Validate{
		logger:  logger,
		parser:  parser,
		tickets: tickets,
	}
}

func (action Validate) Validate(request request.Request, ticket ticket.Ticket) (_ user.User, err error) {
	action.logger.TryToValidate(request, ticket.Nonce())

	ticketUser, nonce, err := action.parser.Parse(ticket.Token())
	if err != nil {
		action.logger.FailedToValidate(request, ticket.Nonce(), err)
		return
	}

	if nonce != ticket.Nonce() {
		err = errValidateDifferentNonce
		action.logger.FailedToValidateBecauseDifferentInfo(request, ticket.Nonce(), err)
		return
	}

	dataUser, expires, found, err := action.tickets.FindUserAndExpires(nonce)
	if err != nil {
		action.logger.FailedToValidate(request, ticket.Nonce(), err)
		return
	}
	if !found {
		err = errValidateNotFoundTicket
		action.logger.FailedToValidateBecauseDifferentInfo(request, ticket.Nonce(), err)
		return
	}
	if ticketUser.ID() != dataUser.ID() {
		err = errValidateDifferentUser
		action.logger.FailedToValidateBecauseDifferentInfo(request, ticket.Nonce(), err)
		return
	}

	if request.RequestedAt().Expired(expires) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateBecauseExpired(request, ticket.Nonce(), err)
		return
	}

	action.logger.AuthByTicket(request, dataUser, ticket.Nonce())
	return dataUser, nil
}
