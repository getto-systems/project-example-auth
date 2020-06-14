package serializer

import (
	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/http_handler/auth_handler/token"

	"github.com/getto-systems/project-example-id/basic"

	"time"
)

type AwsCloudFrontSerializer struct {
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	resource   string
	keyPairID  token.AwsCloudFrontKeyPairID
}

func NewAwsCloudFrontSerializer(pem []byte, resource string, keyPairID token.AwsCloudFrontKeyPairID) AwsCloudFrontSerializer {
	return AwsCloudFrontSerializer{
		privateKey: pem,
		resource:   resource,
		keyPairID:  keyPairID,
	}
}

func (serializer AwsCloudFrontSerializer) Token(ticket basic.Ticket) (token.AwsCloudFrontToken, error) {
	signature, err := serializer.privateKey.Sign(serializer.resource, time.Time(ticket.Expires))
	if err != nil {
		return token.AwsCloudFrontToken{}, err
	}

	return token.AwsCloudFrontToken{
		Policy:    token.AwsCloudFrontPolicy(signature.Policy),
		Signature: token.AwsCloudFrontSignature(signature.Signature),
		KeyPairID: serializer.keyPairID,
	}, nil
}
