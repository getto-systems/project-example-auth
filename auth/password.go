package auth

import (
	"time"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type PasswordAuthenticator interface {
	UserFactory() user.UserFactory
	UserPasswordFactory() user.UserPasswordFactory
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
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

	ticket, err := user.NewTicket(now, param.Path)
	if err != nil {
		return nil, ErrUserAccessDenied
	}

	ticketSerializer := authenticator.TicketSerializer()

	ticketToken, err := ticketSerializer.Token(ticket)
	if err != nil {
		return nil, ErrTicketTokenSerializeFailed
	}

	awsCloudFrontToken, err := authenticator.AwsCloudFrontSerializer().Token(ticket)
	if err != nil {
		return nil, ErrAwsCloudFrontTokenSerializeFailed
	}

	handler(ticket, Token{
		TicketToken:        ticketToken,
		AwsCloudFrontToken: awsCloudFrontToken,
	})

	info, err := ticketSerializer.Info(ticket)
	if err != nil {
		return nil, ErrTicketInfoSerializeFailed
	}

	return info, nil
}
