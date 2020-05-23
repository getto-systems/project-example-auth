package auth

import (
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

func Password(authenticator PasswordAuthenticator, param PasswordParam, handler TokenHandler) (token.TicketInfo, error) {
	userPassword := authenticator.UserPasswordFactory().NewUserPassword(param.UserID)

	if !userPassword.Match(param.Password) {
		return nil, ErrUserPasswordDidNotMatch
	}

	now := authenticator.Now()

	user := authenticator.UserFactory().NewUser(param.UserID)

	ticket, err := user.NewTicket(param.Path, now)
	if err != nil {
		return nil, ErrUserAccessDenied
	}

	err = handleTicketToken(authenticator, ticket, handler)
	if err != nil {
		return nil, err
	}

	info, err := authenticator.TicketSerializer().Info(ticket)
	if err != nil {
		return nil, ErrTicketInfoSerializeFailed
	}

	return info, nil
}
