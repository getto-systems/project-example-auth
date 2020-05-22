package cloudfront

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/getto-systems/project-example-id/signature"
)

type Signer struct {
	pem       Pem
	baseURL   BaseURL
	keyPairID signature.CloudFrontKeyPairID
}

type Pem []byte
type BaseURL string

func NewSigner(pem Pem, baseURL BaseURL, keyPairID signature.CloudFrontKeyPairID) *Signer {
	return &Signer{
		pem,
		baseURL,
		keyPairID,
	}
}

func (signer *Signer) Sign(expires time.Time) (*signature.CloudFrontToken, error) {
	block, _ := pem.Decode(signer.pem)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rng := rand.Reader

	policy := []byte(fmt.Sprintf(
		"{\"Statement\":[{\"Resource\":\"%s/*\",\"Condition\":{\"DateLessThan\":{\"AWS:EpochTime\":%d}}}]}",
		signer.baseURL,
		expires.Unix(),
	))

	hashed := sha1.Sum(policy)

	signed, err := rsa.SignPKCS1v15(rng, key, crypto.SHA1, hashed[:])
	if err != nil {
		return nil, err
	}

	return &signature.CloudFrontToken{
		KeyPairID: signer.keyPairID,
		Policy:    signature.CloudFrontPolicy(cloudfrontBase64(policy)),
		Signature: signature.CloudFrontSignature(cloudfrontBase64(signed)),
	}, nil
}

func cloudfrontBase64(raw []byte) string {
	encoded := base64.StdEncoding.EncodeToString(raw)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "=", "_")
	encoded = strings.ReplaceAll(encoded, "/", "~")
	return encoded
}
