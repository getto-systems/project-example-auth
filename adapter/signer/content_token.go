package signer

import (
	"time"

	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type ContentTokenSigner struct {
	privateKey aws_cloudfront_token.KeyPairPrivateKey
	resource   string
}

func NewContentTokenSigner(pem []byte, resource string) ContentTokenSigner {
	return ContentTokenSigner{
		privateKey: pem,
		resource:   resource,
	}
}

func (signer ContentTokenSigner) signer() ticket.ContentTokenSigner {
	return signer
}

func (signer ContentTokenSigner) Sign(expires data.Expires) (_ ticket.ContentToken, err error) {
	signedToken, err := signer.privateKey.Sign(signer.resource, time.Time(expires))
	if err != nil {
		return
	}

	return ContentToken{signedToken: signedToken}, nil
}

type ContentToken struct {
	signedToken aws_cloudfront_token.SignedCookieToken
}

func (token ContentToken) contentToken() ticket.ContentToken {
	return token
}

func (token ContentToken) Policy() string {
	return token.signedToken.Policy
}

func (token ContentToken) Signature() string {
	return token.signedToken.Signature
}
