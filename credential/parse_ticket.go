package credential

import (
	infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errValidateMatchFailedNonce = data.NewError("Ticket.Validate", "MatchFailed.Nonce")
	errValidateNotFoundTicket   = data.NewError("Ticket.Validate", "NotFound.Ticket")
	errValidateMatchFailedUser  = data.NewError("Ticket.Validate", "MatchFailed.User")
	errValidateAlreadyExpired   = data.NewError("Ticket.Validate", "AlreadyExpired")
)

type ParseTicket struct {
	logger infra.ParseTicketLogger
	parser infra.TicketParser
}

func NewParseTicket(logger infra.ParseTicketLogger, parser infra.TicketParser) ParseTicket {
	return ParseTicket{
		logger: logger,
		parser: parser,
	}
}

func (action ParseTicket) Parse(request request.Request, ticket credential.Ticket) (_ user.User, err error) {
	action.logger.TryToParseTicket(request, ticket.Nonce())

	user, nonce, err := action.parser.Parse(ticket.Signature())
	if err != nil {
		action.logger.FailedToParseTicket(request, ticket.Nonce(), err)
		return
	}
	if nonce != ticket.Nonce() {
		err = errValidateMatchFailedNonce
		action.logger.FailedToParseTicketBecauseNonceMatchFailed(request, ticket.Nonce(), err)
		return
	}

	action.logger.ParseTicket(request, ticket.Nonce(), user)
	return user, nil
}