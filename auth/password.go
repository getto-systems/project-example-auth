package auth

import (
	"github.com/getto-systems/project-example-id/basic"

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
		"PasswordParam{RequestedAt:%s, UserID:%s, Password:[MASKED], Path:%s}",
		param.RequestedAt.String(),
		param.UserID,
		param.Path,
	)
}

func Password(authenticator PasswordAuthenticator, param PasswordParam) (basic.Ticket, error) {
	logger := authenticator.Logger()

	logger.Debugf("password auth: %v", param)

	userPassword := authenticator.UserPasswordFactory().New(param.UserID)

	password, err := userPassword.Password()
	if err != nil {
		logger.Auditf("user password not found: %s; %s", err, param.UserID)
		return basic.Ticket{}, ErrUserPasswordNotFound
	}

	err = password.Match(param.Password)
	if err != nil {
		logger.Auditf("password match failed: %s; %s", err, param.UserID)
		return basic.Ticket{}, ErrUserPasswordDidNotMatch
	}

	user := authenticator.UserFactory().New(param.UserID)

	logger.Debugf("new ticket: %v", param)

	ticket, err := user.NewTicket(param.Path, param.RequestedAt)
	if err != nil {
		logger.Auditf("access denied: %s; %v", err, param)
		return basic.Ticket{}, ErrUserAccessDenied
	}

	return ticket, nil
}
