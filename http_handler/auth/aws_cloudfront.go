package auth

import (
	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/token"
)

func SetAwsCloudFrontCookie(setter http_handler.CookieSetter, token token.AwsCloudFrontToken) {
	setter.SetCookie(&http_handler.Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(token.KeyPairID),
	})

	setter.SetCookie(&http_handler.Cookie{
		Name:  "CloudFront-Policy",
		Value: string(token.Policy),
	})

	setter.SetCookie(&http_handler.Cookie{
		Name:  "CloudFront-Signature",
		Value: string(token.Signature),
	})
}
