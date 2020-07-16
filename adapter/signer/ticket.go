package signer

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type TicketSigner struct {
	jwt JWTSigner
}

func NewTicketSigner(jwt JWTSigner) TicketSigner {
	return TicketSigner{
		jwt: jwt,
	}
}

func (signer TicketSigner) Signer() ticket.Signer {
	return signer
}

func (signer TicketSigner) Verify(ticket ticket.Ticket) (ticket.Nonce, data.User, data.Expires, error) {
	claims, err := signer.jwt.Parse(string(ticket))
	if err != nil {
		return "", data.User{}, data.Expires{}, err
	}

	nonce := parseNonce(claims["jti"])
	user := parseUser(claims["sub"])
	expires := parseExpires(claims["exp"])

	return nonce, user, expires, nil
}
func parseNonce(raw interface{}) ticket.Nonce {
	nonce, ok := raw.(string)
	if !ok {
		return ""
	}

	return ticket.Nonce(nonce)
}
func parseUser(raw interface{}) data.User {
	userID, ok := raw.(string)
	if !ok {
		return data.User{}
	}

	return data.NewUser(data.UserID(userID))
}
func parseExpires(raw interface{}) data.Expires {
	timeString, ok := raw.(string)
	if !ok {
		return data.Expires{}
	}

	unix, err := strconv.Atoi(timeString)
	if err != nil {
		return data.Expires{}
	}

	return data.Expires(time.Unix(int64(unix), 0))
}

func (signer TicketSigner) Sign(nonce ticket.Nonce, user data.User, expires data.Expires) (ticket.Ticket, error) {
	token, err := signer.jwt.Sign(jwt.MapClaims{
		"sub": user.UserID(),
		"exp": strconv.Itoa(int(time.Time(expires).Unix())),
		"jti": nonce,
	})
	if err != nil {
		return nil, err
	}

	return ticket.Ticket(token), nil
}
