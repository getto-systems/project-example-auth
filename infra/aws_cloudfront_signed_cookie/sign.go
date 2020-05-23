package aws_cloudfront_signed_cookie

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
)

type Token struct {
	Policy    Policy
	Signature Signature
}

type Pem []byte
type BaseURL string

type Policy []byte
type Signature []byte

func Sign(privateKey Pem, baseURL BaseURL, expires time.Time) (Token, error) {
	var nullToken Token

	policy := cloudfrontPolicy(baseURL, expires)

	block, _ := pem.Decode(privateKey)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nullToken, err
	}

	rng := rand.Reader

	hashed := sha1.Sum(policy)
	signed, err := rsa.SignPKCS1v15(rng, key, crypto.SHA1, hashed[:])
	if err != nil {
		return nullToken, err
	}

	return Token{
		Policy:    Policy(cloudfrontBase64(policy)),
		Signature: Signature(cloudfrontBase64(signed)),
	}, nil
}

func cloudfrontPolicy(baseURL BaseURL, expires time.Time) []byte {
	// AWS CloudFront flavored json : no extra spaces, strict formatted
	return []byte(fmt.Sprintf(
		"{\"Statement\":[{\"Resource\":\"%s\",\"Condition\":{\"DateLessThan\":{\"AWS:EpochTime\":%d}}}]}",
		baseURL,
		expires.Unix(),
	))
}

func cloudfrontBase64(raw []byte) string {
	// AWS CloudFront flavored base64
	encoded := base64.StdEncoding.EncodeToString(raw)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "=", "_")
	encoded = strings.ReplaceAll(encoded, "/", "~")
	return encoded
}
