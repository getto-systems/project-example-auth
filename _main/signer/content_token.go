package signer

import (
	"time"

	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/credential"
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

func (signer ContentTokenSigner) signer() infra.ContentTokenSigner {
	return signer
}

func (signer ContentTokenSigner) Sign(expires credential.TokenExpires) (_ credential.ContentKeyID, _ credential.ContentPolicy, _ credential.ContentSignature, err error) {
	signed, err := signer.privateKey.Sign(signer.resource, time.Time(expires))
	if err != nil {
		return
	}

	return signer.keyID,
		credential.ContentPolicy(signed.Policy),
		credential.ContentSignature(signed.Signature),
		nil
}
