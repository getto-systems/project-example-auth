package auth

import (
	"github.com/getto-systems/project-example-id/token"
)

func (setter CookieSetter) setAwsCloudFrontCookie(token token.AwsCloudFrontToken) {
	setter.setCookie(&Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(token.KeyPairID),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Policy",
		Value: string(token.Policy),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Signature",
		Value: string(token.Signature),
	})
}
