package data

type (
	User struct {
		userID UserID
	}
	UserID string
)

func NewUser(userID UserID) User {
	return User{
		userID: userID,
	}
}

func (user User) UserID() UserID {
	return user.userID
}
