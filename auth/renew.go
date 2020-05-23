package auth

import (
	"fmt"
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

func (param RenewParam) String() string {
	return fmt.Sprintf(
		"RenewParam{TicketToken:%s, Path:%s}",
		param.TicketToken,
		param.Path,
	)
}

func Renew(authenticator RenewAuthenticator, param RenewParam, handler TokenHandler) (token.TicketInfo, error) {
	logger := authenticator.Logger()

	ticketSerializer := authenticator.TicketSerializer()

	logger.Debugf("renew token: %v", param)

	ticket, err := ticketSerializer.Parse(param.TicketToken, param.Path)
	if err != nil {
		logger.Debugf("parse token error: %s; %v", err, param)
		return nil, ErrTicketTokenParseFailed
	}

	now := authenticator.Now()

	logger.Debugf("check renew required: %v/%s", ticket, now)

	if ticket.IsRenewRequired(now) {
		logger.Debugf("renew token: %v/%s", ticket, now)

		user := authenticator.UserFactory().NewUser(ticket.UserID())

		new_ticket, err := user.NewTicket(param.Path, now)
		if err != nil {
			logger.Auditf("access denied: %s; %v; %v", err, ticket, param)
			return nil, ErrUserAccessDenied
		}

		logger.Auditf("token renewed: %v; %s", new_ticket, param.Path)
		ticket = new_ticket

		err = handleTicketToken(authenticator, ticket, handler)
		if err != nil {
			return nil, err
		}
	}

	logger.Debugf("serialize ticket info: %v", ticket)

	info, err := ticketSerializer.Info(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return nil, ErrTicketInfoSerializeFailed
	}

	return info, nil
}
