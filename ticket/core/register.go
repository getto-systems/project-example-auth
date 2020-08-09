package ticket_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Register(request request.Request, user user.User, exp credential.Expiration) (_ credential.TicketNonce, _ credential.Expires, err error) {
	expires := exp.Expires(request)
	limit := exp.ExtendLimit(request)

	action.logger.TryToRegister(request, user, expires, limit)

	nonce, err := action.tickets.RegisterTicket(action.gen, user, expires, exp.ExpireSecond(), limit)
	if err != nil {
		action.logger.FailedToRegister(request, user, expires, limit, err)
		return
	}

	action.logger.Register(request, user, expires, limit, nonce)
	return nonce, expires, nil
}
