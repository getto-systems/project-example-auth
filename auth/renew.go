package auth

import (
	"time"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type RenewAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
	Now() time.Time
}

type RenewParam struct {
	TicketToken token.TicketToken
	Path        user.Path
}

func Renew(authenticator RenewAuthenticator, param RenewParam, handler TokenHandler) (token.TicketInfo, error) {
	ticketSerializer := authenticator.TicketSerializer()

	ticket, err := ticketSerializer.Parse(param.TicketToken, param.Path)
	if err != nil {
		return nil, ErrTicketTokenParseFailed
	}

	now := authenticator.Now()

	if ticket.IsRenewRequired(now) {
		user := authenticator.UserFactory().NewUser(ticket.UserID())

		ticket, err := user.NewTicket(param.Path, now)
		if err != nil {
			return nil, ErrUserAccessDenied
		}

		err = handleTicketToken(authenticator, ticket, handler)
		if err != nil {
			return nil, err
		}
	}

	info, err := ticketSerializer.Info(ticket)
	if err != nil {
		return nil, ErrTicketInfoSerializeFailed
	}

	return info, nil
}
