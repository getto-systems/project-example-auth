package token

import (
	"github.com/getto-systems/project-example-id/user"
)

type AwsCloudFrontSerializer interface {
	Token(user.Ticket) (AwsCloudFrontToken, error)
}

type AwsCloudFrontToken struct {
	KeyPairID AwsCloudFrontKeyPairID
	Policy    AwsCloudFrontPolicy
	Signature AwsCloudFrontSignature
}

type AwsCloudFrontKeyPairID string
type AwsCloudFrontPolicy string
type AwsCloudFrontSignature string
