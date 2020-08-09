package signer

import (
	gotime "time"

	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/time"
)

type ContentTokenSigner struct {
	keyID      credential.ContentKeyID
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	resource   string
}

func NewContentTokenSigner(keyID credential.ContentKeyID, pem []byte, resource string) ContentTokenSigner {
	return ContentTokenSigner{
		keyID:      keyID,
		privateKey: pem,
		resource:   resource,
	}
}

func (signer ContentTokenSigner) signer() credential.ContentTokenSigner {
	return signer
}

func (signer ContentTokenSigner) Sign(expires time.Expires) (_ credential.ContentToken, err error) {
	signed, err := signer.privateKey.Sign(signer.resource, gotime.Time(expires))
	if err != nil {
		return
	}

	return credential.NewContentToken(
		signer.keyID,
		credential.ContentPolicy(signed.Policy),
		credential.ContentSignature(signed.Signature),
	), nil
}
