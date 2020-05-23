package user

type UserPassword struct {
	db UserPasswordRepository

	userID UserID
}

type UserPasswordRepository interface {
	MatchUserPassword(UserID, Password) bool
}

type Password string

func (userPassword UserPassword) Match(password Password) bool {
	return userPassword.db.MatchUserPassword(userPassword.userID, password)
}

type UserPasswordFactory struct {
	db UserPasswordRepository
}

func NewUserPasswordFactory(db UserPasswordRepository) UserPasswordFactory {
	return UserPasswordFactory{db}
}

func (f UserPasswordFactory) NewUserPassword(userID UserID) UserPassword {
	return UserPassword{f.db, userID}
}
