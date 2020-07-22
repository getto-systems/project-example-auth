package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type RegisterEventPublisher interface {
	GetLogin(data.Request, data.User)
	GetLoginFailed(data.Request, data.User, error)

	RegisterPassword(data.Request, data.User)
	RegisterPasswordFailed(data.Request, data.User, error)
	RegisteredPassword(data.Request, data.User)
}

type RegisterDB interface {
	FilterLogin(data.User) ([]Login, error)
	RegisterPassword(data.User, HashedPassword) error
}

type Generator interface {
	GeneratePassword(RawPassword) (HashedPassword, error)
}
