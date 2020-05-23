package auth

import (
	"errors"
	"fmt"

	"github.com/getto-systems/project-example-id/logger"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type TokenHandler func(user.Ticket, Token)

type Token struct {
	TicketToken        token.TicketToken
	AwsCloudFrontToken token.AwsCloudFrontToken
}

func (token Token) String() string {
	return fmt.Sprintf(
		"Token{TicketToken:%s, AwsCloudFrontToken:%s}",
		token.TicketToken,
		token.AwsCloudFrontToken,
	)
}

var ErrTicketTokenParseFailed = errors.New("ticket token parse failed")
var ErrUserPasswordDidNotMatch = errors.New("user password did not match")
var ErrUserAccessDenied = errors.New("user access denied")
var ErrTicketTokenSerializeFailed = errors.New("ticket token serialize failed")
var ErrAwsCloudFrontTokenSerializeFailed = errors.New("aws cloudfront token serialize failed")
var ErrTicketInfoSerializeFailed = errors.New("ticket info serialize failed")

type Authenticator interface {
	Logger() logger.Logger
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
}

func handleTicketToken(authenticator Authenticator, ticket user.Ticket, handler TokenHandler) error {
	logger := authenticator.Logger()

	logger.Debugf("serialize ticket: %v", ticket)

	logger.Debug("by ticket serializer...")

	ticketToken, err := authenticator.TicketSerializer().Token(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return ErrTicketTokenSerializeFailed
	}

	logger.Debug("by aws cloudfront serializer...")

	awsCloudFrontToken, err := authenticator.AwsCloudFrontSerializer().Token(ticket)
	if err != nil {
		logger.Errorf("aws cloudfront serialize error: %s; %v", err, ticket)
		return ErrAwsCloudFrontTokenSerializeFailed
	}

	logger.Debug("handling ticket token...")

	handler(ticket, Token{
		TicketToken:        ticketToken,
		AwsCloudFrontToken: awsCloudFrontToken,
	})

	return nil
}
