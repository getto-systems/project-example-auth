package signer

import (
	"crypto/ecdsa"

	"github.com/dgrijalva/jwt-go"

	"errors"
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
	verifyKey() interface{}
}

type JWT_ES_512_Key struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type JWT_Pem struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewJWT_ES_512_Key(pem JWT_Pem) (JWT_ES_512_Key, error) {
	var key JWT_ES_512_Key

	if pem.PrivateKey != nil {
		privateKey, err := jwt.ParseECPrivateKeyFromPEM(pem.PrivateKey)
		if err != nil {
			return JWT_ES_512_Key{}, err
		}

		key.privateKey = privateKey
	}

	if pem.PublicKey != nil {
		publicKey, err := jwt.ParseECPublicKeyFromPEM(pem.PublicKey)
		if err != nil {
			return JWT_ES_512_Key{}, err
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

func (key JWT_ES_512_Key) verifyKey() interface{} {
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

func (key JWT_HS_512_Key) verifyKey() interface{} {
	return key.key
}

func (signer JWTSigner) Parse(token string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return signer.key.verifyKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
func (signer JWTSigner) Sign(claims jwt.MapClaims) (string, error) {
	key := signer.key
	jwtToken := jwt.NewWithClaims(key.signingMethod(), claims)
	return jwtToken.SignedString(key.signKey())
}
