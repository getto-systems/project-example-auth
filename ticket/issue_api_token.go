package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type IssueApiTokenEventPublisher interface {
	IssueApiToken(data.Request, data.User, data.Roles, data.Expires)
	IssueApiTokenFailed(data.Request, data.User, data.Roles, data.Expires, error)
}

type IssueApiTokenDB interface {
	FindUserRoles(data.User) (data.Roles, error)
}

type ApiTokenSigner interface {
	Sign(data.User, data.Roles, data.Expires) (ApiToken, error)
}
