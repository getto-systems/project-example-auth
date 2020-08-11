package ticket_core

import (
	"github.com/getto-systems/project-example-id/_misc/errors"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errExtendNotFoundNonce = errors.NewError("Ticket.Extend", "NotFound.Nonce")
)

// user が正しいことは確認済みでなければならない
func (action action) Extend(request request.Request, user user.User, nonce credential.TicketNonce) (_ credential.Ticket, err error) {
	action.logger.TryToExtend(request, user)

	limit, found, err := action.tickets.FindExtendLimit(nonce)
	if err != nil {
		action.logger.FailedToExtend(request, user, err)
		return
	}
	if !found {
		// ticket は validated のはず。ticket が存在しないのはプログラムがおかしいのでエラーログ
		err = errExtendNotFoundNonce
		action.logger.FailedToExtend(request, user, err)
		return
	}

	ticketExpires, tokenExpires := limit.Extend(request, action.ticketExpireSecond, action.tokenExpireSecond)

	action.tickets.UpdateExpires(nonce, ticketExpires)

	action.logger.Extend(request, user, ticketExpires, limit)
	return credential.NewTicket(user, nonce, ticketExpires, tokenExpires), nil
}
