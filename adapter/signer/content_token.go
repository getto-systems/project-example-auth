package signer

import (
	"time"

	"github.com/getto-systems/aws_cloudfront_token-go"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
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

func (signer ContentTokenSigner) ContentTokenSigner() ticket.ContentTokenSigner {
	return signer
}

func (signer ContentTokenSigner) Sign(expires data.Expires) (ticket.ContentToken, error) {
	signedToken, err := signer.privateKey.Sign(signer.resource, time.Time(expires))
	if err != nil {
		return nil, err
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
