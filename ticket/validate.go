package ticket

import (
	"github.com/getto-systems/project-example-id/data/errors"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateDifferentNonce = errors.NewError("Ticket.Validate", "DifferentNonce")
	errValidateAlreadyExpired = errors.NewError("Ticket.Validate", "AlreadyExpired")
)

type Validate struct {
	logger ticket.ValidateLogger
	parser ticket.TicketParser
}

func NewValidate(logger ticket.ValidateLogger, parser ticket.TicketParser) Validate {
	return Validate{
		logger: logger,
		parser: parser,
	}
}

func (action Validate) Validate(request request.Request, ticket ticket.Ticket) (_ user.User, err error) {
	action.logger.TryToValidate(request, ticket.Nonce())

	user, nonce, expires, err := action.parser.Parse(ticket.Token())
	if err != nil {
		action.logger.FailedToValidate(request, ticket.Nonce(), err)
		return
	}

	if nonce != ticket.Nonce() {
		err = errValidateDifferentNonce
		action.logger.FailedToValidate(request, ticket.Nonce(), err)
		return
	}

	if request.RequestedAt().Expired(expires) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateBecauseExpired(request, ticket.Nonce(), err)
		return
	}

	action.logger.AuthByTicket(request, user, ticket.Nonce())
	return user, nil
}
