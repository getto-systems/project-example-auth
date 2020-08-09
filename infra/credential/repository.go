package credential

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	ApiUserRepository interface {
		FindApiRoles(user.User) (credential.ApiRoles, bool, error)

		RegisterApiRoles(user.User, credential.ApiRoles) error
	}
)
