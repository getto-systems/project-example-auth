package auth

import (
	"errors"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type TokenHandler func(user.Ticket, Token)

type Token struct {
	TicketToken        token.TicketToken
	AwsCloudFrontToken token.AwsCloudFrontToken
}

var ErrTicketTokenParseFailed = errors.New("ticket token parse failed")
var ErrUserPasswordDidNotMatch = errors.New("user password did not match")
var ErrUserAccessDenied = errors.New("user access denied")
var ErrTicketTokenSerializeFailed = errors.New("ticket token serialize failed")
var ErrAwsCloudFrontTokenSerializeFailed = errors.New("aws cloudfront token serialize failed")
var ErrTicketInfoSerializeFailed = errors.New("ticket info serialize failed")

type Authenticator interface {
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
}

func handleTicketToken(authenticator Authenticator, ticket user.Ticket, handler TokenHandler) error {
	ticketToken, err := authenticator.TicketSerializer().Token(ticket)
	if err != nil {
		return ErrTicketTokenSerializeFailed
	}

	awsCloudFrontToken, err := authenticator.AwsCloudFrontSerializer().Token(ticket)
	if err != nil {
		return ErrAwsCloudFrontTokenSerializeFailed
	}

	handler(ticket, Token{
		TicketToken:        ticketToken,
		AwsCloudFrontToken: awsCloudFrontToken,
	})

	return nil
}
