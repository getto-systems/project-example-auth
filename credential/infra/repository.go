package infra

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	ApiUserRepository interface {
		FindApiRoles(user.User) (credential.ApiRoles, bool, error)

		RegisterApiRoles(user.User, credential.ApiRoles) error
	}
)
