package ticket_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Register(request request.Request, user user.User, ticketExtendSecond credential.TicketExtendSecond) (_ credential.Ticket, err error) {
	limit := credential.NewTicketExtendLimit(request, ticketExtendSecond)
	expires := credential.NewTicketExpires(request, action.ticketExpireSecond)
	action.logger.TryToRegister(request, user, expires, limit)

	nonce, err := action.tickets.RegisterTicket(action.nonceGenerator, user, expires, limit)
	if err != nil {
		action.logger.FailedToRegister(request, user, expires, limit, err)
		return
	}

	action.logger.Register(request, user, expires, limit)
	return credential.NewTicket(user, nonce, expires, credential.NewTokenExpires(request, action.tokenExpireSecond)), nil
}
