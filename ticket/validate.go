package ticket

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateMatchFailedNonce = data.NewError("Ticket.Validate", "MatchFailed.Nonce")
	errValidateNotFoundTicket   = data.NewError("Ticket.Validate", "NotFound.Ticket")
	errValidateMatchFailedUser  = data.NewError("Ticket.Validate", "MatchFailed.User")
	errValidateAlreadyExpired   = data.NewError("Ticket.Validate", "AlreadyExpired")
)

type Validate struct {
	logger  ticket.ValidateLogger
	tickets ticket.TicketRepository
}

func NewValidate(logger ticket.ValidateLogger, tickets ticket.TicketRepository) Validate {
	return Validate{
		logger:  logger,
		tickets: tickets,
	}
}

// user が正しいことは確認済みでなければならない
func (action Validate) Validate(request request.Request, user user.User, ticket credential.Ticket) (err error) {
	action.logger.TryToValidate(request, user, ticket.Nonce())

	dataUser, expires, found, err := action.tickets.FindUserAndExpires(ticket.Nonce())
	if err != nil {
		action.logger.FailedToValidate(request, user, ticket.Nonce(), err)
		return
	}
	if !found {
		err = errValidateNotFoundTicket
		action.logger.FailedToValidateBecauseTicketNotFound(request, user, ticket.Nonce(), err)
		return
	}
	if user.ID() != dataUser.ID() {
		err = errValidateMatchFailedUser
		action.logger.FailedToValidateBecauseUserMatchFailed(request, user, ticket.Nonce(), err)
		return
	}
	if request.RequestedAt().Expired(expires) {
		err = errValidateAlreadyExpired
		action.logger.FailedToValidateBecauseExpired(request, user, ticket.Nonce(), err)
		return
	}

	action.logger.AuthByTicket(request, dataUser, ticket.Nonce())
	return nil
}
