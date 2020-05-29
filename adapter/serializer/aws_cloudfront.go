package serializer

import (
	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type AwsCloudFrontSerializer struct {
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	baseURL    string
	keyPairID  token.AwsCloudFrontKeyPairID
}

func NewAwsCloudFrontSerializer(pem []byte, baseURL string, keyPairID token.AwsCloudFrontKeyPairID) AwsCloudFrontSerializer {
	return AwsCloudFrontSerializer{
		privateKey: pem,
		baseURL:    baseURL,
		keyPairID:  keyPairID,
	}
}

func (serializer AwsCloudFrontSerializer) Token(ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	var nullToken token.AwsCloudFrontToken

	signature, err := serializer.privateKey.Sign(serializer.baseURL, ticket.Expires())
	if err != nil {
		return nullToken, err
	}

	return token.AwsCloudFrontToken{
		Policy:    token.AwsCloudFrontPolicy(signature.Policy),
		Signature: token.AwsCloudFrontSignature(signature.Signature),
		KeyPairID: serializer.keyPairID,
	}, nil
}
