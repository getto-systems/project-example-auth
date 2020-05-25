package auth

import (
	"github.com/getto-systems/project-example-id/logger"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"errors"
	"fmt"
)

type TokenHandler func(user.Ticket, Token)

type Token struct {
	RenewToken         token.RenewToken
	AwsCloudFrontToken token.AwsCloudFrontToken
}

func (token Token) String() string {
	return fmt.Sprintf(
		"Token{RenewToken:%s, AwsCloudFrontToken:%s}",
		token.RenewToken,
		token.AwsCloudFrontToken,
	)
}

var ErrRenewTokenParseFailed = errors.New("ticket token parse failed")
var ErrUserPasswordDidNotMatch = errors.New("user password did not match")
var ErrUserAccessDenied = errors.New("user access denied")
var ErrRenewTokenSerializeFailed = errors.New("renew token serialize failed")
var ErrAwsCloudFrontTokenSerializeFailed = errors.New("aws cloudfront token serialize failed")
var ErrAppTokenSerializeFailed = errors.New("app token serialize failed")

type Authenticator interface {
	Logger() logger.Logger
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
}

func handleTicket(authenticator Authenticator, ticket user.Ticket, handler TokenHandler) error {
	logger := authenticator.Logger()

	logger.Debugf("serialize ticket: %v", ticket)

	logger.Debug("serialize renew token...")

	renewToken, err := authenticator.TicketSerializer().RenewToken(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return ErrRenewTokenSerializeFailed
	}

	logger.Debug("serialize aws cloudfront token...")

	awsCloudFrontToken, err := authenticator.AwsCloudFrontSerializer().Token(ticket)
	if err != nil {
		logger.Errorf("aws cloudfront serialize error: %s; %v", err, ticket)
		return ErrAwsCloudFrontTokenSerializeFailed
	}

	logger.Debug("handling ticket token...")

	handler(ticket, Token{
		RenewToken:         renewToken,
		AwsCloudFrontToken: awsCloudFrontToken,
	})

	return nil
}
