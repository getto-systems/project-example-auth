package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type ValidateEventPublisher interface {
	ValidatePassword(data.Request, Login)
	ValidatePasswordFailed(data.Request, Login, error)
	AuthenticatedByPassword(data.Request, Login, data.User)
}

type ValidateDB interface {
	FindPasswordByLogin(Login) (data.User, HashedPassword, error)
}

type Matcher interface {
	MatchPassword(HashedPassword, RawPassword) error
}
