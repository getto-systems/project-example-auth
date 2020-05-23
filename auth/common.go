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
var ErrTicketTokenEncodeFailed = errors.New("ticket token encode failed")
var ErrAwsCloudFrontTokenEncodeFailed = errors.New("aws cloudfront token encode failed")
var ErrInfoEncodeFailed = errors.New("info encode failed")
