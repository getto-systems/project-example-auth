package user

import (
	"github.com/getto-systems/project-example-id/basic"
)

type Password struct {
	enc    PasswordEncrypter
	hashed basic.HashedPassword
}

func (p Password) Match(password basic.Password) error {
	return p.enc.MatchPassword(p.hashed, password)
}

type PasswordEncrypter interface {
	GeneratePassword(basic.Password) (basic.HashedPassword, error)
	MatchPassword(basic.HashedPassword, basic.Password) error
}

type UserPassword struct {
	db  UserPasswordRepository
	enc PasswordEncrypter

	userID basic.UserID
}

type UserPasswordRepository interface {
	UserPassword(basic.UserID) (basic.HashedPassword, error)
}

func (p UserPassword) Password() (Password, error) {
	hashed, err := p.db.UserPassword(p.userID)
	if err != nil {
		return Password{}, err
	}

	return Password{
		enc:    p.enc,
		hashed: hashed,
	}, nil
}

type UserPasswordFactory struct {
	db  UserPasswordRepository
	enc PasswordEncrypter
}

func NewUserPasswordFactory(db UserPasswordRepository, enc PasswordEncrypter) UserPasswordFactory {
	return UserPasswordFactory{
		db:  db,
		enc: enc,
	}
}

func (f UserPasswordFactory) New(userID basic.UserID) UserPassword {
	return UserPassword{
		db:     f.db,
		enc:    f.enc,
		userID: userID,
	}
}
