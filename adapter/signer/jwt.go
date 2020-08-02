package signer

import (
	"crypto/ecdsa"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidJWT = errors.New("invalid jwt")
)

type JWTSigner struct {
	key JWTKey
}

func NewJWTSigner(key JWTKey) JWTSigner {
	return JWTSigner{
		key: key,
	}
}

type JWTKey interface {
	signingMethod() jwt.SigningMethod
	signKey() interface{}
	validateKey() interface{}
}

type JWT_ES_512_Key struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type JWT_Pem struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewJWT_ES_512_Key(pem JWT_Pem) (_ JWT_ES_512_Key, err error) {
	var key JWT_ES_512_Key

	if pem.PrivateKey != nil {
		privateKey, parseErr := jwt.ParseECPrivateKeyFromPEM(pem.PrivateKey)
		if parseErr != nil {
			err = parseErr
			return
		}

		key.privateKey = privateKey
	}

	if pem.PublicKey != nil {
		publicKey, parseErr := jwt.ParseECPublicKeyFromPEM(pem.PublicKey)
		if parseErr != nil {
			err = parseErr
			return
		}

		key.publicKey = publicKey
	}

	return key, nil
}

func (JWT_ES_512_Key) signingMethod() jwt.SigningMethod {
	return jwt.SigningMethodES512
}

func (key JWT_ES_512_Key) signKey() interface{} {
	return key.privateKey
}

func (key JWT_ES_512_Key) validateKey() interface{} {
	return key.publicKey
}

type JWT_HS_512_Key struct {
	key []byte
}

func NewJWT_HS_512_Key(key []byte) JWT_HS_512_Key {
	return JWT_HS_512_Key{key: key}
}

func (JWT_HS_512_Key) signingMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS512
}

func (key JWT_HS_512_Key) signKey() interface{} {
	return key.key
}

func (key JWT_HS_512_Key) validateKey() interface{} {
	return key.key
}

func (signer JWTSigner) Parse(token string) (_ jwt.MapClaims, err error) {
	var claims jwt.MapClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return signer.key.validateKey(), nil
	})
	if err != nil {
		return
	}

	if !jwtToken.Valid {
		err = ErrInvalidJWT
		return
	}

	return claims, nil
}
func (signer JWTSigner) Sign(claims jwt.MapClaims) (string, error) {
	key := signer.key
	jwtToken := jwt.NewWithClaims(key.signingMethod(), claims)
	return jwtToken.SignedString(key.signKey())
}
