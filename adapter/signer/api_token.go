package signer

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type ApiTokenSigner struct {
	jwt JWTSigner
}

func NewApiTokenSigner(jwt JWTSigner) ApiTokenSigner {
	return ApiTokenSigner{
		jwt: jwt,
	}
}

func (signer ApiTokenSigner) signer() ticket.ApiTokenSigner {
	return signer
}

func (signer ApiTokenSigner) Sign(user data.User, roles data.Roles, expires data.Expires) (_ ticket.ApiToken, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.UserID(),
		"aud": roles,
		"exp": time.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return ticket.ApiToken(token), nil
}
