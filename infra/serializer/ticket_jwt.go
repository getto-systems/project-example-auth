package serializer

import (
	"crypto/ecdsa"
	"log"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/basic"

	"github.com/getto-systems/project-example-id/http_handler/auth_handler/token"

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
	var key TicketJWT_ES_512_Key

	if pem.PrivateKey != nil {
		privateKey, err := jwt.ParseECPrivateKeyFromPEM(pem.PrivateKey)
		if err != nil {
			return TicketJWT_ES_512_Key{}, err
		}

		key.privateKey = privateKey
	}

	if pem.PublicKey != nil {
		publicKey, err := jwt.ParseECPublicKeyFromPEM(pem.PublicKey)
		if err != nil {
			return TicketJWT_ES_512_Key{}, err
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

func (serializer TicketJWTSerializer) Parse(raw token.RenewToken, path basic.Path) (basic.Ticket, error) {
	var claims jwt.MapClaims
	jwtToken, err := jwt.ParseWithClaims(string(raw), &claims, func(token *jwt.Token) (interface{}, error) {
		return serializer.RenewKey.verifyKey(), nil
	})
	if err != nil {
		return basic.Ticket{}, err
	}

	if !jwtToken.Valid {
		return basic.Ticket{}, errors.New("invalid renew jwt")
	}

	return basic.Ticket{
		UserID:     parseUserID(claims["sub"]),
		Roles:      parseRoles(claims["aud"]),
		Authorized: basic.RequestedAt(parseTime(claims["iat"])),
		Expires:    basic.Expires(parseTime(claims["exp"])),
	}, nil
}
func parseUserID(raw interface{}) basic.UserID {
	userID, ok := raw.(string)
	if !ok {
		return basic.UserID("")
	}

	return basic.UserID(userID)
}
func parseRoles(raw interface{}) basic.Roles {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	var roles []string

	for _, val := range arr {
		role, ok := val.(string)
		if !ok {
			return nil
		}

		roles = append(roles, role)
	}

	return roles
}
func parseTime(raw interface{}) time.Time {
	unixSecond, ok := raw.(int64)
	if !ok {
		var defaultTime time.Time
		return defaultTime
	}

	return time.Unix(unixSecond, 0)
}

func (serializer TicketJWTSerializer) RenewToken(ticket basic.Ticket) (token.RenewToken, error) {
	key := serializer.RenewKey

	jwtToken := jwt.NewWithClaims(key.signingMethod(), jwt.MapClaims{
		"sub": ticket.UserID,
		"aud": ticket.Roles,
		"iat": time.Time(ticket.Authorized).Unix(),
		"exp": time.Time(ticket.Expires).Unix(),
	})

	tokenString, err := jwtToken.SignedString(key.signKey())
	if err != nil {
		log.Println("sign error : renew")
		return nil, err
	}

	return token.RenewToken(tokenString), nil
}

func (serializer TicketJWTSerializer) AppToken(ticket basic.Ticket) (token.AppToken, error) {
	key := serializer.AppKey

	jwtToken := jwt.NewWithClaims(key.signingMethod(), jwt.MapClaims{
		"sub": ticket.UserID,
		"aud": ticket.Roles,
		"iat": time.Time(ticket.Authorized).Unix(),
		"exp": time.Time(ticket.Expires).Unix(),
	})

	tokenString, err := jwtToken.SignedString(key.signKey())
	if err != nil {
		log.Println("sign error : app")
		return token.AppToken{}, err
	}

	return token.AppToken{
		Token:  tokenString,
		UserID: ticket.UserID,
		Roles:  ticket.Roles,
	}, nil
}
