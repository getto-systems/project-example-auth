package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type ValidateEventPublisher interface {
	ValidatePassword(data.Request, data.User)
	ValidatePasswordFailed(data.Request, data.User, error)
	AuthenticatedByPassword(data.Request, data.User)
}

type ValidateDB interface {
	FindUserPassword(data.User) (HashedPassword, error)
}

type Matcher interface {
	MatchPassword(HashedPassword, RawPassword) error
}
