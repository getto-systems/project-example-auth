package api_token

import (
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	ApiUserRepository interface {
		FindApiRoles(user.User) (ApiRoles, bool, error)

		RegisterApiRoles(user.User, ApiRoles) error
	}
)
