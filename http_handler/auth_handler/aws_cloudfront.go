package auth_handler

import (
	"github.com/getto-systems/project-example-id/basic"

	"fmt"
)

type (
	AwsCloudFrontKeyPairID string
	AwsCloudFrontPolicy    string
	AwsCloudFrontSignature string
)

type AwsCloudFrontToken struct {
	Policy    AwsCloudFrontPolicy
	Signature AwsCloudFrontSignature
}

type AwsCloudFrontTicket struct {
	KeyPairID AwsCloudFrontKeyPairID
	Token     AwsCloudFrontToken
}

type AwsCloudFrontIssuer struct {
	keyPairID  AwsCloudFrontKeyPairID
	serializer AwsCloudFrontSerializer
}

type AwsCloudFrontSerializer interface {
	Serialize(basic.Ticket) (AwsCloudFrontToken, error)
}

func (iss AwsCloudFrontIssuer) Authorized(ticket basic.Ticket) (AwsCloudFrontTicket, error) {
	token, err := iss.serializer.Serialize(ticket)
	if err != nil {
		return AwsCloudFrontTicket{}, err
	}

	return AwsCloudFrontTicket{
		KeyPairID: iss.keyPairID,
		Token:     token,
	}, nil
}

func NewAwsCloudFrontIssuer(keyPairID AwsCloudFrontKeyPairID, serializer AwsCloudFrontSerializer) AwsCloudFrontIssuer {
	return AwsCloudFrontIssuer{
		keyPairID:  keyPairID,
		serializer: serializer,
	}
}

func (token AwsCloudFrontToken) String() string {
	return fmt.Sprintf(
		"Token{Policy:%s, Signature:%s}",
		token.Policy,
		token.Signature,
	)
}
