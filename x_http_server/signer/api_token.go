package signer

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/user"
)

type ApiTokenSigner struct {
	jwt JWTSigner
}

func NewApiTokenSigner(jwt JWTSigner) ApiTokenSigner {
	return ApiTokenSigner{
		jwt: jwt,
	}
}

func (signer ApiTokenSigner) signer() infra.ApiTokenSigner {
	return signer
}

func (signer ApiTokenSigner) Sign(user user.User, roles credential.ApiRoles, expires credential.TokenExpires) (_ credential.ApiSignature, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"aud": roles,
		"exp": time.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return credential.ApiSignature(token), nil
}
