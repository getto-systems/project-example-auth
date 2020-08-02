package signer

import (
	gotime "time"

	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/time"
)

type ContentTokenSigner struct {
	keyID      api_token.ContentKeyID
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	resource   string
}

func NewContentTokenSigner(keyID api_token.ContentKeyID, pem []byte, resource string) ContentTokenSigner {
	return ContentTokenSigner{
		keyID:      keyID,
		privateKey: pem,
		resource:   resource,
	}
}

func (signer ContentTokenSigner) signer() api_token.ContentTokenSigner {
	return signer
}

func (signer ContentTokenSigner) Sign(expires time.Expires) (_ api_token.ContentToken, err error) {
	signed, err := signer.privateKey.Sign(signer.resource, gotime.Time(expires))
	if err != nil {
		return
	}

	return api_token.NewContentToken(
		signer.keyID,
		api_token.ContentPolicy(signed.Policy),
		api_token.ContentSignature(signed.Signature),
	), nil
}
