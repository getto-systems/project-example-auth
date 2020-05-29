package user

type UserPassword struct {
	db  UserPasswordRepository
	enc UserPasswordEncrypter

	userID UserID
}

type UserPasswordRepository interface {
	UserPassword(UserID) HashedPassword
}

type UserPasswordEncrypter interface {
	GenerateUserPassword(Password) (HashedPassword, error)
	MatchUserPassword(HashedPassword, Password) error
}

type (
	Password       string
	HashedPassword []byte
)

func (p UserPassword) Match(password Password) error {
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

func (f UserPasswordFactory) NewUserPassword(userID UserID) UserPassword {
	return UserPassword{
		db:     f.db,
		enc:    f.enc,
		userID: userID,
	}
}
