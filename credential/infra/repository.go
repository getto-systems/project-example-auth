package infra

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/user"
)

type (
	ApiUserRepository interface {
		FindApiRoles(user.User) (credential.ApiRoles, bool, error)

		RegisterApiRoles(user.User, credential.ApiRoles) error
	}
)
