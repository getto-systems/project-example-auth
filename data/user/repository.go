package user

type (
	UserRepository interface {
		FindUser(Login) (User, bool, error)
		FindLogin(User) (Login, bool, error)
		RegisterUser(User, Login) error
	}
)
