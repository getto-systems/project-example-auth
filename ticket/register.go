package ticket

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
	infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type Register struct {
	logger  infra.RegisterLogger
	gen     credential_infra.TicketNonceGenerator
	tickets infra.TicketRepository
}

func NewRegister(logger infra.RegisterLogger, gen credential_infra.TicketNonceGenerator, tickets infra.TicketRepository) Register {
	return Register{
		logger:  logger,
		gen:     gen,
		tickets: tickets,
	}
}

// user が正しいことは確認済みでなければならない
func (action Register) Register(request request.Request, user user.User, exp ticket.Expiration) (_ credential.TicketNonce, _ time.Expires, err error) {
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
