package api_token

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	ApiTokenSigner interface {
		Sign(user.User, ApiRoles, time.Expires) (ApiToken, error)
	}

	ContentTokenSigner interface {
		Sign(time.Expires) (ContentToken, error)
	}
)
