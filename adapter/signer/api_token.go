package signer

import (
	gotime "time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data/api_token"
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

func (signer ApiTokenSigner) signer() api_token.ApiTokenSigner {
	return signer
}

func (signer ApiTokenSigner) Sign(user user.User, roles api_token.ApiRoles, expires time.Expires) (_ api_token.ApiToken, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"aud": roles,
		"exp": gotime.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return api_token.NewApiToken(roles, api_token.ApiSignature(token)), nil
}
