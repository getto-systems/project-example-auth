package token

import (
	"github.com/getto-systems/project-example-id/basic"

	"fmt"
)

type AwsCloudFrontSerializer interface {
	Token(basic.TicketData) (AwsCloudFrontToken, error)
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

type (
	AwsCloudFrontKeyPairID string
	AwsCloudFrontPolicy    string
	AwsCloudFrontSignature string
)
