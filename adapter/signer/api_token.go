package signer

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type ApiTokenSigner struct {
	jwt JWTSigner
}

func NewApiTokenSigner(jwt JWTSigner) ApiTokenSigner {
	return ApiTokenSigner{
		jwt: jwt,
	}
}

func (signer ApiTokenSigner) ApiTokenSigner() ticket.ApiTokenSigner {
	return signer
}

func (signer ApiTokenSigner) Sign(user data.User, roles data.Roles, expires data.Expires) (ticket.ApiToken, error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.UserID(),
		"aud": roles,
		"exp": time.Time(expires).Unix(),
	})
	if err != nil {
		log.Println("app token sign error")
		return nil, err
	}

	return ticket.ApiToken(token), nil
}
