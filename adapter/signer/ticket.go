package signer

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type TicketSigner struct {
	jwt JWTSigner
}

func NewTicketSigner(jwt JWTSigner) TicketSigner {
	return TicketSigner{
		jwt: jwt,
	}
}

func (signer TicketSigner) signer() ticket.Signer {
	return signer
}

func (signer TicketSigner) Parse(ticket ticket.Ticket) (_ ticket.Nonce, _ data.User, _ data.Expires, err error) {
	claims, err := signer.jwt.Parse(string(ticket))
	if err != nil {
		return
	}

	nonce := parseNonce(claims["jti"])
	user := parseUser(claims["sub"])
	expires := parseExpires(claims["exp"])

	return nonce, user, expires, nil
}
func parseNonce(raw interface{}) (_ ticket.Nonce) {
	nonce, ok := raw.(string)
	if !ok {
		return
	}

	return ticket.Nonce(nonce)
}
func parseUser(raw interface{}) (_ data.User) {
	userID, ok := raw.(string)
	if !ok {
		return
	}

	return data.NewUser(data.UserID(userID))
}
func parseExpires(raw interface{}) (_ data.Expires) {
	timeString, ok := raw.(string)
	if !ok {
		return
	}

	unix, err := strconv.Atoi(timeString)
	if err != nil {
		return
	}

	return data.Expires(time.Unix(int64(unix), 0))
}

func (signer TicketSigner) Sign(nonce ticket.Nonce, user data.User, expires data.Expires) (_ ticket.Ticket, err error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.UserID(),
		"exp": strconv.Itoa(int(time.Time(expires).Unix())),
		"jti": nonce,
	})
	if err != nil {
		return
	}

	return ticket.Ticket(token), nil
}
