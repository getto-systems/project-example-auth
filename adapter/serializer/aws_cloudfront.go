package serializer

import (
	"github.com/getto-systems/project-example-id/misc/aws_cloudfront_signed_cookie"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type AwsCloudFrontSerializer struct {
	pem       AwsCloudFrontPem
	baseURL   AwsCloudFrontBaseURL
	keyPairID token.AwsCloudFrontKeyPairID
}

type AwsCloudFrontPem []byte
type AwsCloudFrontBaseURL string

func NewAwsCloudFrontSerializer(pem AwsCloudFrontPem, baseURL AwsCloudFrontBaseURL, keyPairID token.AwsCloudFrontKeyPairID) AwsCloudFrontSerializer {
	return AwsCloudFrontSerializer{
		pem:       pem,
		baseURL:   baseURL,
		keyPairID: keyPairID,
	}
}

func (serializer AwsCloudFrontSerializer) Token(ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	var nullToken token.AwsCloudFrontToken

	signature, err := aws_cloudfront_signed_cookie.Sign(
		aws_cloudfront_signed_cookie.Pem(serializer.pem),
		aws_cloudfront_signed_cookie.BaseURL(serializer.baseURL),
		ticket.Expires(),
	)
	if err != nil {
		return nullToken, err
	}

	return token.AwsCloudFrontToken{
		Policy:    token.AwsCloudFrontPolicy(signature.Policy),
		Signature: token.AwsCloudFrontSignature(signature.Signature),
		KeyPairID: serializer.keyPairID,
	}, nil
}
