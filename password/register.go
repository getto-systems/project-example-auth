package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type RegisterLogger interface {
	TryToGetLogin(data.Request, data.User)
	FailedToGetLogin(data.Request, data.User, error)

	TryToRegister(data.Request, data.User)
	FailedToRegister(data.Request, data.User, error)
	Registered(data.Request, data.User)
}

type RegisterDB interface {
	FilterLogin(data.User) ([]Login, error)
	RegisterPassword(data.User, HashedPassword) error
}

type Generator interface {
	GeneratePassword(RawPassword) (HashedPassword, error)
}
