package ticket

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

var (
	errExtendNotFoundNonce = data.NewError("Ticket.Extend", "NotFound.Nonce")
)

type Extend struct {
	logger  ticket.ExtendLogger
	signer  ticket.TicketSigner
	tickets ticket.TicketRepository
}

func NewExtend(logger ticket.ExtendLogger, signer ticket.TicketSigner, tickets ticket.TicketRepository) Extend {
	return Extend{
		logger:  logger,
		signer:  signer,
		tickets: tickets,
	}
}

func (action Extend) Extend(request request.Request, user user.User, oldTicket ticket.Ticket) (_ ticket.Ticket, _ time.Expires, err error) {
	// user が正しいことは確認済みでなければならない
	action.logger.TryToExtend(request, user, oldTicket.Nonce())

	expireSecond, limit, found, err := action.tickets.FindExpireSecondAndExtendLimit(oldTicket.Nonce())
	if err != nil {
		action.logger.FailedToExtend(request, user, oldTicket.Nonce(), err)
		return
	}
	if !found {
		// この時点で ticket が存在しないのはプログラムがおかしいのでエラーログ
		err = errExtendNotFoundNonce
		action.logger.FailedToExtend(request, user, oldTicket.Nonce(), err)
		return
	}

	expires := request.RequestedAt().Expires(expireSecond).Limit(limit)

	token, err := action.signer.Sign(user, oldTicket.Nonce(), expires)
	if err != nil {
		action.logger.FailedToExtend(request, user, oldTicket.Nonce(), err)
		return
	}

	newTicket := ticket.NewTicket(token, oldTicket.Nonce())

	action.logger.Extend(request, user, newTicket.Nonce(), expires, limit)
	return newTicket, expires, nil
}
