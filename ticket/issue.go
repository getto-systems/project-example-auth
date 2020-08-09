package ticket

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type Issue struct {
	logger  ticket.IssueLogger
	signer  credential.TicketSigner
	gen     credential.TicketNonceGenerator
	tickets ticket.TicketRepository
}

func NewIssue(logger ticket.IssueLogger, signer credential.TicketSigner, gen credential.TicketNonceGenerator, tickets ticket.TicketRepository) Issue {
	return Issue{
		logger:  logger,
		signer:  signer,
		gen:     gen,
		tickets: tickets,
	}
}

func (action Issue) Issue(request request.Request, user user.User, exp ticket.Expiration) (_ credential.Ticket, _ time.Expires, err error) {
	expires := exp.Expires(request)
	limit := exp.ExtendLimit(request)

	action.logger.TryToIssue(request, user, expires, limit)

	nonce, err := action.tickets.RegisterTicket(action.gen, user, expires, exp.ExpireSecond(), limit)
	if err != nil {
		action.logger.FailedToIssue(request, user, expires, limit, err)
		return
	}

	token, err := action.signer.Sign(user, nonce, expires)
	if err != nil {
		action.logger.FailedToIssue(request, user, expires, limit, err)
		return
	}

	ticket := credential.NewTicket(token, nonce)

	action.logger.Issue(request, user, expires, limit, ticket.Nonce())
	return ticket, expires, nil
}
