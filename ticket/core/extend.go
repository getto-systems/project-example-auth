package ticket_core

import (
	"github.com/getto-systems/project-example-id/misc/errors"
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	errExtendNotFoundNonce = errors.NewError("Ticket.Extend", "NotFound.Nonce")
)

// user が正しいことは確認済みでなければならない
func (action action) Extend(request request.Request, user user.User, ticket credential.Ticket) (_ expiration.Expires, err error) {
	action.logger.TryToExtend(request, user, ticket.Nonce())

	limit, found, err := action.tickets.FindExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToExtend(request, user, ticket.Nonce(), err)
		return
	}
	if !found {
		// ticket は validated のはず。ticket が存在しないのはプログラムがおかしいのでエラーログ
		err = errExtendNotFoundNonce
		action.logger.FailedToExtend(request, user, ticket.Nonce(), err)
		return
	}

	expires := request.Extend(limit, action.expireSecond)
	action.tickets.UpdateExpires(ticket.Nonce(), expires)

	action.logger.Extend(request, user, ticket.Nonce(), expires, limit)
	return expires, nil
}
