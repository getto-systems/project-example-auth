package serializer

import (
	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/http_handler/auth_handler"

	"github.com/getto-systems/project-example-id/data"

	"time"
)

type AwsCloudFrontSerializer struct {
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	resource   string
}

func NewAwsCloudFrontSerializer(pem []byte, resource string) AwsCloudFrontSerializer {
	return AwsCloudFrontSerializer{
		privateKey: pem,
		resource:   resource,
	}
}

func (serializer AwsCloudFrontSerializer) Serialize(ticket data.Ticket) (auth_handler.AwsCloudFrontToken, error) {
	signature, err := serializer.privateKey.Sign(serializer.resource, time.Time(ticket.Expires))
	if err != nil {
		return auth_handler.AwsCloudFrontToken{}, err
	}

	return auth_handler.AwsCloudFrontToken{
		Policy:    auth_handler.AwsCloudFrontPolicy(signature.Policy),
		Signature: auth_handler.AwsCloudFrontSignature(signature.Signature),
	}, nil
}
