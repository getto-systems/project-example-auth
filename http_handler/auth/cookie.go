package auth

import (
	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/auth"
)

func SetAuthTokenCookie(setter http_handler.CookieSetter, token auth.Token) {
	SetTicketCookie(setter, token.TicketToken)
	SetAwsCloudFrontCookie(setter, token.AwsCloudFrontToken)
}
