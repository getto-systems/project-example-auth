package auth

import (
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"fmt"
	"time"
)

type RenewAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
	Now() time.Time
}

type RenewParam struct {
	RenewToken token.RenewToken
	Path       user.Path
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

	var nullToken token.AppToken

	ticketSerializer := authenticator.TicketSerializer()

	ticket, err := ticketSerializer.Parse(param.RenewToken, param.Path)
	if err != nil {
		logger.Debugf("parse token error: %s; %v", err, param)
		return nullToken, ErrRenewTokenParseFailed
	}

	now := authenticator.Now()

	logger.Debugf("check renew required: %v/%s", ticket, now)

	if ticket.IsRenewRequired(now) {
		logger.Debugf("renew token: %v/%s", ticket, now)

		user := authenticator.UserFactory().NewUser(ticket.UserID())

		new_ticket, err := user.NewTicket(param.Path, now)
		if err != nil {
			logger.Auditf("access denied: %s; %v; %v", err, ticket, param)
			return nullToken, ErrUserAccessDenied
		}

		logger.Auditf("token renewed: %v; %s", new_ticket, param.Path)

		err = handleTicket(authenticator, new_ticket, handler)
		if err != nil {
			return nullToken, err
		}

		ticket = new_ticket
	}

	logger.Debugf("serialize app token: %v", ticket)

	appToken, err := ticketSerializer.AppToken(ticket)
	if err != nil {
		logger.Errorf("app token serialize error: %s; %v", err, ticket)
		return nullToken, ErrAppTokenSerializeFailed
	}

	return appToken, nil
}
