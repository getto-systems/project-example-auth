package user

type UserPassword string

type UserPasswordRepository interface {
	MatchUserPassword(UserID, UserPassword) bool
}
