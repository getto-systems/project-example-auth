package token

import (
	"fmt"

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

func (token AwsCloudFrontToken) String() string {
	return fmt.Sprintf(
		"Token{KeyPairID:%s, Policy:%s, Signature:%s}",
		token.KeyPairID,
		token.Policy,
		token.Signature,
	)
}

type AwsCloudFrontKeyPairID string
type AwsCloudFrontPolicy string
type AwsCloudFrontSignature string
