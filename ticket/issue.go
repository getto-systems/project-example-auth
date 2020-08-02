package ticket

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type Issue struct {
	logger  ticket.IssueLogger
	exp     expiration
	gen     ticket.NonceGenerator
	signer  ticket.TicketSigner
	tickets ticket.TicketRepository
}

func NewIssue(logger ticket.IssueLogger, signer ticket.TicketSigner, exp ticket.ExpirationParam, gen ticket.NonceGenerator, tickets ticket.TicketRepository) Issue {
	return Issue{
		logger:  logger,
		exp:     newExpiration(exp),
		gen:     gen,
		signer:  signer,
		tickets: tickets,
	}
}

func (action Issue) Issue(request request.Request, user user.User) (_ ticket.Ticket, _ time.Expires, err error) {
	expires := action.exp.Expires(request)
	limit := action.exp.ExtendLimit(request)

	action.logger.TryToIssue(request, user, expires, limit)

	nonce, err := action.tickets.RegisterTicket(action.gen, user, expires, limit)
	if err != nil {
		action.logger.FailedToIssue(request, user, expires, limit, err)
		return
	}

	token, err := action.signer.Sign(user, nonce, expires)
	if err != nil {
		action.logger.FailedToIssue(request, user, expires, limit, err)
		return
	}

	ticket := ticket.NewTicket(token, nonce)

	action.logger.Issue(request, user, expires, limit, ticket.Nonce())
	return ticket, expires, nil
}
