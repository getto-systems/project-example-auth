package auth

import (
	"github.com/getto-systems/project-example-id/basic"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"fmt"
)

type PasswordAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
	UserPasswordFactory() user.UserPasswordFactory
}

type PasswordParam struct {
	RequestedAt basic.RequestedAt

	UserID   basic.UserID
	Password basic.Password
	Path     basic.Path
}

func (param PasswordParam) String() string {
	return fmt.Sprintf(
		"PasswordParam{UserID:%s, Password:[MASKED], Path:%s}",
		param.UserID,
		param.Path,
	)
}

func Password(authenticator PasswordAuthenticator, param PasswordParam, handler TokenHandler) (token.AppToken, error) {
	logger := authenticator.Logger()

	logger.Debugf("password auth: %v", param)

	userPassword := authenticator.UserPasswordFactory().New(param.UserID)

	password, err := userPassword.Password()
	if err != nil {
		logger.Auditf("user password not found: %s; %s", err, param.UserID)
		return token.AppToken{}, ErrUserPasswordNotFound
	}

	err = password.Match(param.Password)
	if err != nil {
		logger.Auditf("password match failed: %s; %s", err, param.UserID)
		return token.AppToken{}, ErrUserPasswordDidNotMatch
	}

	user := authenticator.UserFactory().New(param.UserID)

	logger.Debugf("new ticket: %v", param)

	ticket, err := user.NewTicket(param.Path, param.RequestedAt)
	if err != nil {
		logger.Auditf("access denied: %s; %v", err, param)
		return token.AppToken{}, ErrUserAccessDenied
	}

	err = handleTicket(authenticator, ticket, handler)
	if err != nil {
		return token.AppToken{}, err
	}

	logger.Debugf("serialize app token: %v", ticket)

	appToken, err := authenticator.TicketSerializer().AppToken(ticket.Data())
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return token.AppToken{}, ErrAppTokenSerializeFailed
	}

	return appToken, nil
}
