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
