package ticket

import (
	infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errExtendNotFoundNonce = data.NewError("Ticket.Extend", "NotFound.Nonce")
)

type Extend struct {
	logger  infra.ExtendLogger
	tickets infra.TicketRepository
}

func NewExtend(logger infra.ExtendLogger, tickets infra.TicketRepository) Extend {
	return Extend{
		logger:  logger,
		tickets: tickets,
	}
}

// user が正しいことは確認済みでなければならない
func (action Extend) Extend(request request.Request, user user.User, ticket credential.Ticket) (_ time.Expires, err error) {
	action.logger.TryToExtend(request, user, ticket.Nonce())

	expireSecond, limit, found, err := action.tickets.FindExpireSecondAndExtendLimit(ticket.Nonce())
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

	expires := request.RequestedAt().Expires(expireSecond).Limit(limit)
	action.tickets.UpdateExpires(ticket.Nonce(), expires)

	action.logger.Extend(request, user, ticket.Nonce(), expires, limit)
	return expires, nil
}
