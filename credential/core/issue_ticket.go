package credential_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
)

func (action action) IssueTicket(request request.Request, ticket credential.Ticket) (_ credential.TicketToken, err error) {
	action.logger.TryToIssueTicket(request, ticket.User(), ticket.TicketExpires())

	signature, err := action.ticketSigner.Sign(ticket.User(), ticket.Nonce(), ticket.TicketExpires())
	if err != nil {
		action.logger.FailedToIssueTicket(request, ticket.User(), ticket.TicketExpires(), err)
		return
	}

	action.logger.IssueTicket(request, ticket.User(), ticket.TicketExpires())
	return ticket.NewTicketToken(signature), nil
}
