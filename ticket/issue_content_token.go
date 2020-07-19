package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type IssueContentTokenEventPublisher interface {
	IssueContentToken(data.Request, data.User, data.Expires)
	IssueContentTokenFailed(data.Request, data.User, data.Expires, error)
}

type ContentTokenSigner interface {
	Sign(data.Expires) (ContentToken, error)
}
