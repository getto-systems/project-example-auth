package password_encrypter

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/getto-systems/project-example-auth/password/infra"

	"github.com/getto-systems/project-example-auth/password"
)

var (
	errPasswordEmpty   = errors.New("empty password is not allowed")
	errPasswordTooLong = errors.New("password too long")
)

type Encrypter struct {
	cost int
}

func NewEncrypter(cost int) Encrypter {
	return Encrypter{
		cost: cost,
	}
}

func (enc Encrypter) enc() infra.PasswordEncrypter {
	return enc
}

func (enc Encrypter) GeneratePassword(password password.RawPassword) (_ password.HashedPassword, err error) {
	p, err := newRaw(password)
	if err != nil {
		return
	}

	return p.generate(enc.cost)
}

func (enc Encrypter) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) (_ bool, err error) {
	p, err := newRaw(raw)
	if err != nil {
		return
	}

	err = p.compare(hashed)
	if err != nil {
		return false, nil
	}

	return true, nil
}

type raw []byte

func newRaw(password password.RawPassword) (_ raw, err error) {
	bytes := []byte(password)

	if len(bytes) == 0 {
		err = errPasswordEmpty
		return
	}

	if len(bytes) > 72 {
		err = errPasswordTooLong
		return
	}

	return bytes, nil
}

func (password raw) generate(cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (password raw) compare(hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, password)
}
