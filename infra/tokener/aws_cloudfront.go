package tokener

import (
	"github.com/getto-systems/project-example-id/infra/aws_cloudfront_signed_cookie"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type AwsCloudFrontTokener struct {
	pem       AwsCloudFrontPem
	baseURL   AwsCloudFrontBaseURL
	keyPairID token.AwsCloudFrontKeyPairID
}

type AwsCloudFrontPem []byte
type AwsCloudFrontBaseURL string

func NewAwsCloudFrontTokener(pem AwsCloudFrontPem, baseURL AwsCloudFrontBaseURL, keyPairID token.AwsCloudFrontKeyPairID) AwsCloudFrontTokener {
	return AwsCloudFrontTokener{
		pem,
		baseURL,
		keyPairID,
	}
}

func (tokener AwsCloudFrontTokener) Token(ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	var nullToken token.AwsCloudFrontToken

	signature, err := aws_cloudfront_signed_cookie.Sign(
		aws_cloudfront_signed_cookie.Pem(tokener.pem),
		aws_cloudfront_signed_cookie.BaseURL(tokener.baseURL),
		ticket.Expires(),
	)
	if err != nil {
		return nullToken, err
	}

	return token.AwsCloudFrontToken{
		Policy:    token.AwsCloudFrontPolicy(signature.Policy),
		Signature: token.AwsCloudFrontSignature(signature.Signature),
		KeyPairID: tokener.keyPairID,
	}, nil
}
