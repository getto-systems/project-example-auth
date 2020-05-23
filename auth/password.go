package auth

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type PasswordAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
	UserPasswordFactory() user.UserPasswordFactory
	Now() time.Time
}

type PasswordParam struct {
	UserID   user.UserID
	Password user.Password
	Path     user.Path
}

func (param PasswordParam) String() string {
	return fmt.Sprintf(
		"PasswordParam{UserID:%s, Password:[MASKED], Path:%s}",
		param.UserID,
		param.Path,
	)
}

func Password(authenticator PasswordAuthenticator, param PasswordParam, handler TokenHandler) (token.TicketInfo, error) {
	logger := authenticator.Logger()

	logger.Debugf("password auth: %v", param)

	userPassword := authenticator.UserPasswordFactory().NewUserPassword(param.UserID)

	if !userPassword.Match(param.Password) {
		logger.Auditf("password did not match: %s", param.UserID)
		return nil, ErrUserPasswordDidNotMatch
	}

	now := authenticator.Now()

	user := authenticator.UserFactory().NewUser(param.UserID)

	logger.Debugf("new ticket: %v", param)

	ticket, err := user.NewTicket(param.Path, now)
	if err != nil {
		logger.Auditf("access denied: %s; %v", err, param)
		return nil, ErrUserAccessDenied
	}

	err = handleTicketToken(authenticator, ticket, handler)
	if err != nil {
		return nil, err
	}

	logger.Debugf("serialize ticket info: %v", ticket)

	info, err := authenticator.TicketSerializer().Info(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return nil, ErrTicketInfoSerializeFailed
	}

	return info, nil
}
