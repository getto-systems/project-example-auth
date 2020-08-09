package ticket_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/ticket"
)

// user が正しいことは確認済みでなければならない
func (action action) Register(request request.Request, user user.User, exp ticket.Expiration) (_ credential.TicketNonce, _ time.Expires, err error) {
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
