package signer

import (
	"strconv"
	gotime "time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data/ticket"
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

func (signer TicketSigner) sign() ticket.TicketSign {
	return signer
}

func (signer TicketSigner) Parse(token ticket.Token) (_ user.User, _ ticket.Nonce, _ time.Expires, err error) {
	claims, err := signer.jwt.Parse(string(token))
	if err != nil {
		return
	}

	nonce := parseNonce(claims["jti"])
	user := parseUser(claims["sub"])
	expires := parseExpires(claims["exp"])

	return user, nonce, expires, nil
}
func parseNonce(raw interface{}) (_ ticket.Nonce) {
	nonce, ok := raw.(string)
	if !ok {
		return
	}

	return ticket.Nonce(nonce)
}
func parseUser(raw interface{}) (_ user.User) {
	userID, ok := raw.(string)
	if !ok {
		return
	}

	return user.NewUser(user.UserID(userID))
}
func parseExpires(raw interface{}) (_ time.Expires) {
	timeString, ok := raw.(string)
	if !ok {
		return
	}

	unix, err := strconv.Atoi(timeString)
	if err != nil {
		return
	}

	return time.Expires(gotime.Unix(int64(unix), 0))
}

func (signer TicketSigner) Sign(user user.User, nonce ticket.Nonce, expires time.Expires) (_ ticket.Token, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.ID(),
		"exp": strconv.Itoa(int(gotime.Time(expires).Unix())),
		"jti": nonce,
	})
	if err != nil {
		return
	}

	return ticket.Token(token), nil
}
