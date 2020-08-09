package signer

import (
	gotime "time"

	"github.com/dgrijalva/jwt-go"

	infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
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

func (signer ApiTokenSigner) Sign(user user.User, roles credential.ApiRoles, expires time.Expires) (_ credential.ApiToken, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"aud": roles,
		"exp": gotime.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return credential.NewApiToken(roles, credential.ApiSignature(token)), nil
}
