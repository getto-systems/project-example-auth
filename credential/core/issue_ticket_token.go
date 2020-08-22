package credential_core

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
)

func (action action) IssueTicketToken(request request.Request, ticket credential.Ticket) (_ credential.TicketToken, err error) {
	action.logger.TryToIssueTicketToken(request, ticket.User(), ticket.TicketExpires())

	signature, err := action.ticketSigner.Sign(ticket.User(), ticket.Nonce(), ticket.TicketExpires())
	if err != nil {
		action.logger.FailedToIssueTicketToken(request, ticket.User(), ticket.TicketExpires(), err)
		return
	}

	action.logger.IssueTicketToken(request, ticket.User(), ticket.TicketExpires())
	return ticket.NewTicketToken(signature), nil
}
