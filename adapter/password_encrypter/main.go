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

type PasswordEncrypter struct {
	cost int
}

func NewPasswordEncrypter(cost int) PasswordEncrypter {
	return PasswordEncrypter{
		cost: cost,
	}
}

func (enc PasswordEncrypter) matcher() password.Matcher {
	return enc
}

func (enc PasswordEncrypter) gen() password.Generator {
	return enc
}

func (enc PasswordEncrypter) GeneratePassword(password password.RawPassword) (password.HashedPassword, error) {
	p, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	return p.generate(enc.cost)
}

func (enc PasswordEncrypter) MatchPassword(hashed password.HashedPassword, password password.RawPassword) error {
	p, err := NewPassword(password)
	if err != nil {
		return err
	}

	err = p.compare(hashed)
	if err != nil {
		return err
	}

	return nil
}

type Password []byte

func NewPassword(password password.RawPassword) (Password, error) {
	bytes := []byte(password)

	if len(bytes) == 0 {
		return nil, ErrPasswordEmpty
	}

	if len(bytes) > 72 {
		return nil, ErrPasswordTooLong
	}

	return bytes, nil
}

func (password Password) generate(cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (password Password) compare(hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, password)
}
