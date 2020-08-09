package signer

import (
	"strconv"
	gotime "time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type TicketSigner struct {
	jwt JWTSigner
}

func NewTicketSigner(jwt JWTSigner) TicketSigner {
	return TicketSigner{
		jwt: jwt,
	}
}

func (signer TicketSigner) sign() api_token.TicketSign {
	return signer
}

func (signer TicketSigner) Parse(signature api_token.TicketSignature) (_ user.User, _ api_token.TicketNonce, err error) {
	claims, err := signer.jwt.Parse(string(signature))
	if err != nil {
		return
	}

	nonce := parseNonce(claims["jti"])
	user := parseUser(claims["sub"])

	return user, nonce, nil
}
func parseNonce(raw interface{}) (_ api_token.TicketNonce) {
	nonce, ok := raw.(string)
	if !ok {
		return
	}

	return api_token.TicketNonce(nonce)
}
func parseUser(raw interface{}) (_ user.User) {
	userID, ok := raw.(string)
	if !ok {
		return
	}

	return user.NewUser(user.UserID(userID))
}

func (signer TicketSigner) Sign(user user.User, nonce api_token.TicketNonce, expires time.Expires) (_ api_token.TicketSignature, err error) {
	signature, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"exp": strconv.Itoa(int(gotime.Time(expires).Unix())),
		"jti": nonce,
	})
	if err != nil {
		return
	}

	return api_token.TicketSignature(signature), nil
}
