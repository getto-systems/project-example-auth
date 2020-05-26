package serializer

import (
	"crypto/ecdsa"
	"log"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"errors"
	"time"
)

type TicketJWTSerializer struct {
	RenewKey TicketJWTKey
	AppKey   TicketJWTKey
}

type TicketJWTKey interface {
	signingMethod() jwt.SigningMethod
	signKey() interface{}
	verifyKey() interface{}
}

type TicketJWT_ES_512_Key struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type TicketJWT_ES_512_Pem struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewTicketJWT_ES_512_Key(pem TicketJWT_ES_512_Pem) (TicketJWT_ES_512_Key, error) {
	var nullKey TicketJWT_ES_512_Key

	var key TicketJWT_ES_512_Key

	if pem.PrivateKey != nil {
		privateKey, err := jwt.ParseECPrivateKeyFromPEM(pem.PrivateKey)
		if err != nil {
			return nullKey, err
		}

		key.privateKey = privateKey
	}

	if pem.PublicKey != nil {
		publicKey, err := jwt.ParseECPublicKeyFromPEM(pem.PublicKey)
		if err != nil {
			return nullKey, err
		}

		key.publicKey = publicKey
	}

	return key, nil
}

func (TicketJWT_ES_512_Key) signingMethod() jwt.SigningMethod {
	return jwt.SigningMethodES512
}

func (key TicketJWT_ES_512_Key) signKey() interface{} {
	return key.privateKey
}

func (key TicketJWT_ES_512_Key) verifyKey() interface{} {
	return key.publicKey
}

type TicketJWT_HS_512_Key struct {
	key []byte
}

func NewTicketJWT_HS_512_Key(key []byte) TicketJWT_HS_512_Key {
	return TicketJWT_HS_512_Key{key: key}
}

func (TicketJWT_HS_512_Key) signingMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS512
}

func (key TicketJWT_HS_512_Key) signKey() interface{} {
	return key.key
}

func (key TicketJWT_HS_512_Key) verifyKey() interface{} {
	return key.key
}

func (serializer TicketJWTSerializer) Parse(raw token.RenewToken, path user.Path) (user.Ticket, error) {
	var nullTicket user.Ticket

	var claims jwt.MapClaims
	jwtToken, err := jwt.ParseWithClaims(string(raw), &claims, func(token *jwt.Token) (interface{}, error) {
		return serializer.RenewKey.verifyKey(), nil
	})
	if err != nil {
		return nullTicket, err
	}

	if !jwtToken.Valid {
		return nullTicket, errors.New("invalid renew jwt")
	}

	return user.RestrictTicket(path, user.TicketData{
		UserID:     parseUserID(claims["sub"]),
		Roles:      parseRoles(claims["aud"]),
		Authorized: parseTime(claims["iat"]),
		Expires:    parseTime(claims["exp"]),
	})
}
func parseUserID(raw interface{}) user.UserID {
	var nullUserID user.UserID

	userID, ok := raw.(string)
	if !ok {
		return nullUserID
	}

	return user.UserID(userID)
}
func parseRoles(raw interface{}) user.Roles {
	roles, ok := raw.([]string)
	if !ok {
		return nil
	}

	return roles
}
func parseTime(raw interface{}) time.Time {
	var nullTime time.Time

	unixSecond, ok := raw.(int64)
	if !ok {
		return nullTime
	}

	return time.Unix(unixSecond, 0)
}

func (serializer TicketJWTSerializer) RenewToken(ticket user.Ticket) (token.RenewToken, error) {
	key := serializer.RenewKey

	jwtToken := jwt.NewWithClaims(key.signingMethod(), jwt.MapClaims{
		"sub": ticket.UserID(),
		"aud": ticket.Roles(),
		"iat": ticket.Authorized().Unix(),
		"exp": ticket.Expires().Unix(),
	})

	tokenString, err := jwtToken.SignedString(key.signKey())
	if err != nil {
		log.Println("sign error : renew")
		return nil, err
	}

	return token.RenewToken(tokenString), nil
}

func (serializer TicketJWTSerializer) AppToken(ticket user.Ticket) (token.AppToken, error) {
	var nullToken token.AppToken

	key := serializer.AppKey

	jwtToken := jwt.NewWithClaims(key.signingMethod(), jwt.MapClaims{
		"sub": ticket.UserID(),
		"aud": ticket.Roles(),
		"iat": ticket.Authorized().Unix(),
		"exp": ticket.Expires().Unix(),
	})

	tokenString, err := jwtToken.SignedString(key.signKey())
	if err != nil {
		log.Println("sign error : app")
		return nullToken, err
	}

	return token.AppToken{
		Token:  tokenString,
		UserID: ticket.UserID(),
		Roles:  ticket.Roles(),
	}, nil
}
