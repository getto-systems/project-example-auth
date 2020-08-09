package credential

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	Logger interface {
		IssueApiTokenLogger
		IssueContentTokenLogger
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, time.Expires)
		FailedToIssueApiToken(request.Request, user.User, time.Expires, error)
		IssueApiToken(request.Request, user.User, ApiRoles, time.Expires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, time.Expires)
		FailedToIssueContentToken(request.Request, user.User, time.Expires, error)
		IssueContentToken(request.Request, user.User, time.Expires)
	}
)
