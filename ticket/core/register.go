package ticket_core

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Register(request request.Request, user user.User, second expiration.ExtendSecond) (_ credential.TicketNonce, _ expiration.Expires, err error) {
	limit := request.NewExtendLimit(second)
	expires := request.NewExpires(action.expireSecond)

	action.logger.TryToRegister(request, user, expires, limit)

	nonce, err := action.tickets.RegisterTicket(action.nonceGenerator, user, expires, limit)
	if err != nil {
		action.logger.FailedToRegister(request, user, expires, limit, err)
		return
	}

	action.logger.Register(request, user, expires, limit)
	return nonce, expires, nil
}
