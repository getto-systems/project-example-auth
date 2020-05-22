package signature

import (
	"time"
)

type CloudFrontSigner interface {
	Sign(time.Time) (*CloudFrontToken, error)
}

type CloudFrontToken struct {
	KeyPairID CloudFrontKeyPairID
	Policy    CloudFrontPolicy
	Signature CloudFrontSignature
}

type CloudFrontKeyPairID string
type CloudFrontPolicy string
type CloudFrontSignature string
