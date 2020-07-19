package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type RegisterEventPublisher interface {
	RegisterPassword(data.Request, data.User)
	RegisterPasswordFailed(data.Request, data.User, error)
	PasswordRegistered(data.Request, data.User)
}

type RegisterDB interface {
	RegisterUserPassword(data.User, HashedPassword) error
}

type Generator interface {
	GeneratePassword(RawPassword) (HashedPassword, error)
}
