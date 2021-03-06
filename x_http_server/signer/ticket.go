package signer

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/user"
)

type TicketSigner struct {
	jwt JWTSigner
}

func NewTicketSigner(jwt JWTSigner) TicketSigner {
	return TicketSigner{
		jwt: jwt,
	}
}

func (signer TicketSigner) sign() infra.TicketSign {
	return signer
}

func (signer TicketSigner) Parse(signature credential.TicketSignature) (_ user.User, _ credential.TicketNonce, err error) {
	claims, err := signer.jwt.Parse(string(signature))
	if err != nil {
		return
	}

	nonce := parseNonce(claims["jti"])
	user := parseUser(claims["sub"])

	return user, nonce, nil
}
func parseNonce(raw interface{}) (_ credential.TicketNonce) {
	nonce, ok := raw.(string)
	if !ok {
		return
	}

	return credential.TicketNonce(nonce)
}
func parseUser(raw interface{}) (_ user.User) {
	userID, ok := raw.(string)
	if !ok {
		return
	}

	return user.NewUser(user.UserID(userID))
}

func (signer TicketSigner) Sign(user user.User, nonce credential.TicketNonce, expires credential.TicketExpires) (_ credential.TicketSignature, err error) {
	signature, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"exp": strconv.Itoa(int(time.Time(expires).Unix())),
		"jti": nonce,
	})
	if err != nil {
		return
	}

	return credential.TicketSignature(signature), nil
}
