package auth

import (
	"github.com/getto-systems/project-example-id/basic"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"fmt"
)

type RenewAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
}

type RenewParam struct {
	RequestedAt basic.Time
	RenewToken  token.RenewToken
	Path        basic.Path
}

func (param RenewParam) String() string {
	return fmt.Sprintf(
		"RenewParam{RenewParam:%s, Path:%s}",
		param.RenewToken,
		param.Path,
	)
}

func Renew(authenticator RenewAuthenticator, param RenewParam, handler TokenHandler) (token.AppToken, error) {
	logger := authenticator.Logger()

	logger.Debugf("renew token: %v", param)

	ticketSerializer := authenticator.TicketSerializer()

	ticket, err := ticketSerializer.Parse(param.RenewToken, param.Path)
	if err != nil {
		logger.Debugf("parse token error: %s; %v", err, param)
		return token.AppToken{}, ErrRenewTokenParseFailed
	}

	logger.Debugf("check renew required: %v/%s", ticket, param.RequestedAt)

	if ticket.IsRenewRequired(param.RequestedAt) {
		logger.Debugf("renew token: %v/%s", ticket, param.RequestedAt)

		user := authenticator.UserFactory().NewUser(ticket.UserID())

		new_ticket, err := user.NewTicket(param.Path, param.RequestedAt)
		if err != nil {
			logger.Auditf("access denied: %s; %v; %v", err, ticket, param)
			return token.AppToken{}, ErrUserAccessDenied
		}

		logger.Auditf("token renewed: %v; %s", new_ticket, param.Path)

		err = handleTicket(authenticator, new_ticket, handler)
		if err != nil {
			return token.AppToken{}, err
		}

		ticket = new_ticket
	}

	logger.Debugf("serialize app token: %v", ticket)

	appToken, err := ticketSerializer.AppToken(ticket)
	if err != nil {
		logger.Errorf("app token serialize error: %s; %v", err, ticket)
		return token.AppToken{}, ErrAppTokenSerializeFailed
	}

	return appToken, nil
}
