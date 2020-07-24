package password_encrypter

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/getto-systems/project-example-id/password"

	"errors"
)

var (
	ErrPasswordEmpty   = errors.New("empty password is not allowed")
	ErrPasswordTooLong = errors.New("password too long")
)

type Encrypter struct {
	cost int
}

func NewEncrypter(cost int) Encrypter {
	return Encrypter{
		cost: cost,
	}
}

func (enc Encrypter) matcher() password.PasswordMatcher {
	return enc
}

func (enc Encrypter) gen() password.PasswordGenerator {
	return enc
}

func (enc Encrypter) GeneratePassword(password password.RawPassword) (_ password.HashedPassword, err error) {
	p, err := NewPassword(password)
	if err != nil {
		return
	}

	return p.generate(enc.cost)
}

func (enc Encrypter) MatchPassword(hashed password.HashedPassword, password password.RawPassword) (err error) {
	p, err := NewPassword(password)
	if err != nil {
		return
	}

	err = p.compare(hashed)
	if err != nil {
		return
	}

	return nil
}

type Password []byte

func NewPassword(password password.RawPassword) (_ Password, err error) {
	bytes := []byte(password)

	if len(bytes) == 0 {
		err = ErrPasswordEmpty
		return
	}

	if len(bytes) > 72 {
		err = ErrPasswordTooLong
		return
	}

	return bytes, nil
}

func (password Password) generate(cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (password Password) compare(hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, password)
}
