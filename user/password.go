package user

import (
	"github.com/getto-systems/project-example-id/basic"
)

type UserPassword struct {
	db  UserPasswordRepository
	enc UserPasswordEncrypter

	userID basic.UserID
}

type UserPasswordRepository interface {
	UserPassword(basic.UserID) basic.HashedPassword
}

type UserPasswordEncrypter interface {
	GenerateUserPassword(basic.Password) (basic.HashedPassword, error)
	MatchUserPassword(basic.HashedPassword, basic.Password) error
}

func (p UserPassword) Match(password basic.Password) error {
	hashed := p.db.UserPassword(p.userID)
	return p.enc.MatchUserPassword(hashed, password)
}

type UserPasswordFactory struct {
	db  UserPasswordRepository
	enc UserPasswordEncrypter
}

func NewUserPasswordFactory(db UserPasswordRepository, enc UserPasswordEncrypter) UserPasswordFactory {
	return UserPasswordFactory{
		db:  db,
		enc: enc,
	}
}

func (f UserPasswordFactory) NewUserPassword(userID basic.UserID) UserPassword {
	return UserPassword{
		db:     f.db,
		enc:    f.enc,
		userID: userID,
	}
}
