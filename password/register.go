package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type RegisterEventPublisher interface {
	GetLogin(data.Request, data.User)
	LoginNotFound(data.Request, data.User, error)

	RegisterPassword(data.Request, data.User)
	RegisterPasswordFailed(data.Request, data.User, error)
	PasswordRegistered(data.Request, data.User)
}

type RegisterDB interface {
	FindLoginByUser(data.User) (Login, error)
	RegisterPasswordOfUser(data.User, HashedPassword) error
}

type Generator interface {
	GeneratePassword(RawPassword) (HashedPassword, error)
}
